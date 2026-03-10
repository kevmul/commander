package workflow

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ─── Command ──────────────────────────────────────────────────────────────────

func (e *Executor) executeCommand(step Step, stepNum, totalSteps int) error {
	command := e.parser.Parse(step.Command)

	if step.Description != "" {
		desc := e.parser.Parse(step.Description)
		fmt.Printf("[%d/%d] %s\n", stepNum, totalSteps, desc)
	} else {
		fmt.Printf("[%d/%d] Running: %s\n", stepNum, totalSteps, command)
	}

	if step.Interactive {
		cmd := exec.Command("sh", "-c", command)
		cmd.Env = e.env.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil && !step.IgnoreError {
			return fmt.Errorf("command failed: %w", err)
		}
		return nil
	}

	if step.CaptureEnv {
		cmd := exec.Command("sh", "-c", command)
		cmd.Env = e.env.Environ()
		cmd.Stdin = os.Stdin

		var buf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &buf)

		err := cmd.Run()
		if err != nil && !step.IgnoreError {
			return fmt.Errorf("command failed: %w", err)
		}

		// Parse stdout/stderr for KEY=VALUE pairs and store in workflow env.
		// Also push into the template parser so they're available as {{KEY}}
		// in subsequent steps.
		e.env.ParseAndApply(buf.String())
		for key, value := range e.env.vars {
			e.parser.Set(key, value)
		}

		return nil
	}

	if step.CaptureOutput {
		cmd := exec.Command("sh", "-c", command)
		cmd.Env = e.env.Environ()

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		output := strings.TrimSpace(stdout.String())
		if step.OutputVariable != "" {
			e.parser.Set(step.OutputVariable, output)
		}

		if err != nil {
			if !step.IgnoreError {
				return fmt.Errorf("command failed: %w\nStderr: %s", err, stderr.String())
			}
			fmt.Printf("⚠️  Command failed but continuing: %v\n", err)
		}
		return nil
	}

	// Default: stream output directly, no capture
	cmd := exec.Command("sh", "-c", command)
	cmd.Env = e.env.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil && !step.IgnoreError {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
