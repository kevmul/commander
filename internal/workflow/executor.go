package workflow

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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

func (e *Executor) executeInput(step Step) error {
	prompt := e.parser.Parse(step.Prompt)

	cmd := exec.Command("gum", "input", "--placeholder", prompt)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get input: %w", err)
	}

	value := strings.TrimSpace(string(output))
	e.parser.Set(step.Variable, value)

	return nil
}

func (e *Executor) executeSelect(step Step) error {
	prompt := e.parser.Parse(step.Prompt)

	args := []string{"choose"}
	args = append(args, step.Options...)
	args = append(args, "--header", prompt)

	cmd := exec.Command("gum", args...)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get selection: %w", err)
	}

	value := strings.TrimSpace(string(output))
	e.parser.Set(step.Variable, value)

	return nil
}

func (e *Executor) executeConfirm(step Step) error {
	prompt := e.parser.Parse(step.Prompt)

	cmd := exec.Command("gum", "confirm", prompt)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		// User said no
		e.parser.Set(step.Variable, "false")
	} else {
		// User said yes
		e.parser.Set(step.Variable, "true")
	}

	return nil
}

func (e *Executor) executeCommand(step Step, stepNum, totalSteps int) error {
	command := e.parser.Parse(step.Command)

	if step.Description != "" {
		desc := e.parser.Parse(step.Description)
		fmt.Printf("[%d/%d] %s\n", stepNum, totalSteps, desc)
	} else {
		fmt.Printf("[%d/%d] Running: %s\n", stepNum, totalSteps, command)
	}

	cmd := exec.Command("sh", "-c", command)

	if step.CaptureOutput {
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		output := strings.TrimSpace(stdout.String())
		if step.OutputVariable != "" {
			e.parser.Set(step.OutputVariable, output)
		}

		if err != nil {
			if step.DieOnError {
				return fmt.Errorf("command failed: %w\nStderr: %s", err, stderr.String())
			}
			fmt.Printf("⚠️  Command failed but continuing: %v\n", err)
		}
	} else {
		var stdout bytes.Buffer
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil && step.DieOnError {
			return fmt.Errorf("command failed: %w", err)
		}

		output := strings.TrimSpace(stdout.String())
		if step.OutputVariable != "" {
			e.parser.Set(step.OutputVariable, output)
		}

	}

	return nil
}
