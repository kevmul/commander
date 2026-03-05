package workflow

import (
	"fmt"
	"github.com/kevmul/cmdr/internal/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ─── Input ────────────────────────────────────────────────────────────────────

type inputModel struct {
	prompt   string
	helpText string
	value    string
	done     bool
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
	cursor := styles.CursorStyle.Render("‣")
	return lipgloss.JoinVertical(lipgloss.Left,
		m.prompt,
		styles.HelpTextStyle.Render(m.helpText),
		styles.InputStyle.Render(fmt.Sprintf("%s %s", cursor, m.value)),
	)
}

func (e *Executor) executeInput(step Step) error {
	prompt := e.parser.Parse(step.Prompt)
	helpText := e.parser.Parse(step.HelpText)

	m := inputModel{helpText: helpText, prompt: fmt.Sprintf("%s:", prompt)}
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
