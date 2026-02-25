package workflow

import (
	"fmt"

	"github.com/kevmul/cmdr/internal/template"
)

// Executor runs workflows
type Executor struct {
	parser *template.Parser
}

// NewExecutor creates a new workflow executor
func NewExecutor() *Executor {
	return &Executor{
		parser: template.NewParser(),
	}
}

// Execute runs a workflow
func (e *Executor) Execute(workflow *Workflow) error {
	e.parser.Reset()

	fmt.Printf("\nRunning workflow: %s\n", workflow.Name)
	if workflow.Description != "" {
		fmt.Printf("   %s\n", workflow.Description)
	}
	fmt.Println()

	for i, step := range workflow.Steps {
		if err := e.executeStep(step, i+1, len(workflow.Steps)); err != nil {
			return err
		}
	}

	fmt.Println("\n✅ Workflow completed successfully!")
	return nil
}

func (e *Executor) executeStep(step Step, stepNum, totalSteps int) error {
	switch step.Type {
	case StepTypeInput:
		return e.executeInput(step)
	case StepTypeSelect:
		return e.executeSelect(step)
	case StepTypeConfirm:
		return e.executeConfirm(step)
	case StepTypeCommand:
		return e.executeCommand(step, stepNum, totalSteps)
	default:
		return fmt.Errorf("unknown step type: %s", step.Type)
	}
}
