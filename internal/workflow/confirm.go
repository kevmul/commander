package workflow

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ─── Confirm ──────────────────────────────────────────────────────────────────

type confirmModel struct {
	prompt    string
	confirmed bool
	done      bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "y":
			m.confirmed = true
			m.done = true
			return m, tea.Quit
		case "n":
			m.confirmed = false
			m.done = true
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	return fmt.Sprintf("%s (y/n): ", m.prompt)
}

func (e *Executor) executeConfirm(step Step) error {
	prompt := e.parser.Parse(step.Prompt)

	m := confirmModel{prompt: prompt}
	p := tea.NewProgram(m)

	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("confirm failed: %w", err)
	}

	final := result.(confirmModel)
	if !final.done {
		return fmt.Errorf("confirm cancelled")
	}

	answer := "false"
	label := "No"
	if final.confirmed {
		answer = "true"
		label = "Yes"
	}
	fmt.Printf("  ✔ %s\n\n", label)

	e.parser.Set(step.Variable, answer)
	return nil
}

// ─── Command ──────────────────────────────────────────────────────────────────

func (e *Executor) executeCommand(step Step, stepNum, totalSteps int) error {
	command := e.parser.Parse(step.Command)

	if step.Description != "" {
		desc := e.parser.Parse(step.Description)
		fmt.Printf("[%d/%d] %s\n", stepNum, totalSteps, desc)
	} else {
		fmt.Printf("[%d/%d] Running: %s\n", stepNum, totalSteps, command)
	}

	cmd := exec.Command("sh", "-c", command)

	if step.Interactive {
		cmd.Stdin = os.Stdin
	}

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
