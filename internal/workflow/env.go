package workflow

import (
	"os"
	"strings"
	"sync"
)

// WorkflowEnv holds env vars captured during workflow execution.
// Workflow-captured vars override the user's shell environment.
// It is safe for concurrent use.
type WorkflowEnv struct {
	mu   sync.RWMutex
	vars map[string]string
}

// NewWorkflowEnv creates a new WorkflowEnv
func NewWorkflowEnv() *WorkflowEnv {
	return &WorkflowEnv{vars: make(map[string]string)}
}

// Set stores a key/value pair in the workflow env store.
// If the key already exists it is overwritten.
func (e *WorkflowEnv) Set(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.vars[key] = value
}

// Get retrieves a value by key.
func (e *WorkflowEnv) Get(key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.vars[key]
	return v, ok
}

// Reset clears all captured env vars.
func (e *WorkflowEnv) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.vars = make(map[string]string)
}

// Environ returns a merged []string slice suitable for exec.Cmd.Env.
// It starts with os.Environ() and workflow-captured vars are merged on
// top, overriding any shell-level values with the same key.
func (e *WorkflowEnv) Environ() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Index workflow vars for O(1) lookup during merge
	overrides := make(map[string]string, len(e.vars))
	for k, v := range e.vars {
		overrides[k] = v
	}

	base := os.Environ()
	result := make([]string, 0, len(base)+len(overrides))

	// Add shell env, skipping keys that workflow overrides
	for _, entry := range base {
		key, _, _ := strings.Cut(entry, "=")
		if _, overridden := overrides[key]; !overridden {
			result = append(result, entry)
		}
	}

	// Append workflow vars
	for k, v := range overrides {
		result = append(result, k+"="+v)
	}

	return result
}

// ParseAndApply parses a block of command output for shell export lines
// and stores any discovered key/value pairs.
//
// Supported formats:
//
//	export KEY=VALUE
//	export KEY="VALUE WITH SPACES"
//	export KEY='VALUE WITH SPACES'
//	KEY=VALUE
//
// Lines that do not match are silently skipped, making this safe to
// run against any command output regardless of content.
func (e *WorkflowEnv) ParseAndApply(output string) {
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		// If line contains spaces before KEY=VALUE, take only the last space-separated token
		if idx := strings.LastIndex(line, " "); idx != -1 && strings.Contains(line[idx:], "=") {
			line = line[idx+1:]
		}

		line = strings.TrimSpace(line)

		// Strip leading `export ` if present
		line = strings.TrimPrefix(line, "export ")

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Key must be non-empty and a valid identifier (no spaces)
		if key == "" || strings.ContainsAny(key, " \t") {
			continue
		}

		// Strip surrounding quotes from value
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		e.Set(key, value)
	}
}
