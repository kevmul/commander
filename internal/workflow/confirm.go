package workflow

import (
	"fmt"
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
