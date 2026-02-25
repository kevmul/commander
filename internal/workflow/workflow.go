package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// StepType represents the type of step in a workflow
type StepType string

const (
	StepTypeInput   StepType = "input"
	StepTypeSelect  StepType = "select"
	StepTypeConfirm StepType = "confirm"
	StepTypeCommand StepType = "command"
)

// Step represents a single step in a workflow
type Step struct {
	Type           StepType `json:"type"`
	Prompt         string   `json:"prompt,omitempty"`
	HelpText       string   `json:"helpText,omitempty"`
	Variable       string   `json:"variable,omitempty"`
	Options        []string `json:"options,omitempty"`
	Command        string   `json:"command,omitempty"`
	Description    string   `json:"description,omitempty"`
	CaptureOutput  bool     `json:"capture_output,omitempty"`
	OutputVariable string   `json:"output_variable,omitempty"`
	DieOnError     bool     `json:"die_on_error,omitempty"`
}

// Workflow represents a complete workflow with multiple steps
type Workflow struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Steps       []Step `json:"steps"`
}

var (
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)
	leadingTrailing = regexp.MustCompile(`^-+|-+$`)
)

// Slugify converts a string into a lowercase, hyphen-separated slug.
// e.g. "My Cool Workflow!" -> "my-cool-workflow"
func Slugify(s string) string {
	slug := strings.ToLower(s)
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	slug = leadingTrailing.ReplaceAllString(slug, "")
	return slug
}

// Store handles loading and saving workflows
type Store struct {
	filePath string
}

// NewStore creates a new workflow store
func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "cmdr")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &Store{filePath: filepath.Join(configDir, "workflows.json")}, nil
}

// readAll reads all workflows from the file
func (s *Store) readAll() ([]Workflow, error) {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return []Workflow{}, nil
	}
	if err != nil {
		return nil, err
	}

	var workflows []Workflow
	if err := json.Unmarshal(data, &workflows); err != nil {
		return nil, err
	}
	return workflows, nil
}

// writeAll writes all workflows to the file
func (s *Store) writeAll(workflows []Workflow) error {
	data, err := json.MarshalIndent(workflows, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

// List returns all available workflows
func (s *Store) List() ([]Workflow, error) {
	return s.readAll()
}

// Load loads a workflow by key
func (s *Store) Load(key string) (*Workflow, error) {
	workflows, err := s.readAll()
	if err != nil {
		return nil, err
	}

	for _, w := range workflows {
		if w.Key == key {
			return &w, nil
		}
	}
	return nil, fmt.Errorf("workflow not found: %s", key)
}

// Save saves a workflow, matching on key for updates.
// The workflow's Key must be set before calling Save.
func (s *Store) Save(workflow *Workflow) error {
	if workflow.Key == "" {
		workflow.Key = Slugify(workflow.Name)
	}

	workflows, err := s.readAll()
	if err != nil {
		return err
	}

	for i, w := range workflows {
		if w.Key == workflow.Key {
			workflows[i] = *workflow
			return s.writeAll(workflows)
		}
	}

	workflows = append(workflows, *workflow)
	return s.writeAll(workflows)
}

// Delete deletes a workflow by key
func (s *Store) Delete(key string) error {
	workflows, err := s.readAll()
	if err != nil {
		return err
	}

	for i, w := range workflows {
		if w.Key == key {
			workflows = append(workflows[:i], workflows[i+1:]...)
			return s.writeAll(workflows)
		}
	}
	return fmt.Errorf("workflow not found: %s", key)
}

// KeyExists checks whether a key is already in use
func (s *Store) KeyExists(key string) bool {
	workflows, err := s.readAll()
	if err != nil {
		return false
	}

	for _, w := range workflows {
		if w.Key == key {
			return true
		}
	}
	return false
}

// Exists checks if a workflow exists by key (kept for backwards compat)
func (s *Store) Exists(key string) bool {
	return s.KeyExists(key)
}
