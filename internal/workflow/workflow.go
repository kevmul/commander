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
	configDir string
}

// NewStore creates a new workflow store
func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "cmdr")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &Store{configDir: configDir}, nil
}

// List returns all available workflows
func (s *Store) List() ([]Workflow, error) {
	entries, err := os.ReadDir(s.configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var workflows []Workflow
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		workflow, err := s.Load(entry.Name()[:len(entry.Name())-5]) // Remove .json
		if err != nil {
			continue // Skip invalid workflows
		}
		workflows = append(workflows, *workflow)
	}

	return workflows, nil
}

// Load loads a workflow by name
func (s *Store) Load(name string) (*Workflow, error) {
	path := filepath.Join(s.configDir, name+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow Workflow
	if err := json.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow: %w", err)
	}

	return &workflow, nil
}

// Save saves a workflow
func (s *Store) Save(workflow *Workflow) error {
	path := filepath.Join(s.configDir, workflow.Name+".json")

	data, err := json.MarshalIndent(workflow, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal workflow: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	return nil
}

// Delete deletes a workflow by name
func (s *Store) Delete(name string) error {
	path := filepath.Join(s.configDir, name+".json")

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete workflow: %w", err)
	}

	return nil
}

// Exists checks if a workflow exists
func (s *Store) Exists(name string) bool {
	path := filepath.Join(s.configDir, name+".json")
	_, err := os.Stat(path)
	return err == nil
}
