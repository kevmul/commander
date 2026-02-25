package workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Steps       []Step `json:"steps"`
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

// Load loads a workflow by name
func (s *Store) Load(name string) (*Workflow, error) {
	workflows, err := s.readAll()
	if err != nil {
		return nil, err
	}

	for _, w := range workflows {
		if w.Name == name {
			return &w, nil
		}
	}
	return nil, fmt.Errorf("workflow not found: %s", name)
}

// Save saves a workflow
func (s *Store) Save(workflow *Workflow) error {
	workflows, err := s.readAll()
	if err != nil {
		return err
	}

	for i, w := range workflows {
		if w.Name == workflow.Name {
			workflows[i] = *workflow
			return s.writeAll(workflows)
		}
	}

	workflows = append(workflows, *workflow)
	return s.writeAll(workflows)
}

// Delete deletes a workflow by name
func (s *Store) Delete(name string) error {
	workflows, err := s.readAll()
	if err != nil {
		return err
	}

	for i, w := range workflows {
		if w.Name == name {
			workflows = append(workflows[:i], workflows[i+1:]...)
			return s.writeAll(workflows)
		}
	}
	return fmt.Errorf("workflow not found: %s", name)
}

// Exists checks if a workflow exists
func (s *Store) Exists(name string) bool {
	workflows, err := s.readAll()
	if err != nil {
		return false
	}

	for _, w := range workflows {
		if w.Name == name {
			return true
		}
	}
	return false
}
