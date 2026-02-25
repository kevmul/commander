package workflow

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ─── Input ────────────────────────────────────────────────────────────────────

type inputModel struct {
	prompt string
	value  string
	done   bool
}

func (m inputModel) Init() tea.Cmd { return nil }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.done = true
			return m, tea.Quit
		case tea.KeyBackspace, tea.KeyDelete:
			if len(m.value) > 0 {
				m.value = m.value[:len(m.value)-1]
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			m.done = false
			return m, tea.Quit
		default:
			if msg.Type == tea.KeyRunes {
				m.value += string(msg.Runes)
			}
		}
	}
	return m, nil
}

func (m inputModel) View() string {
	return fmt.Sprintf("%s %s▌", m.prompt, m.value)
}

func (e *Executor) executeInput(step Step) error {
	prompt := e.parser.Parse(step.Prompt)

	m := inputModel{prompt: fmt.Sprintf("? %s:", prompt)}
	p := tea.NewProgram(m)

	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("input failed: %w", err)
	}

	final := result.(inputModel)
	if !final.done {
		return fmt.Errorf("input cancelled")
	}

	fmt.Println() // newline after inline input
	e.parser.Set(step.Variable, final.value)
	return nil
}
