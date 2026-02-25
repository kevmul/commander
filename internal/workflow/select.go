package workflow

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ─── Select ───────────────────────────────────────────────────────────────────

type selectModel struct {
	prompt   string
	options  []string
	cursor   int
	selected string
	done     bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			m.selected = m.options[m.cursor]
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("? %s:\n", m.prompt))
	for i, opt := range m.options {
		if i == m.cursor {
			sb.WriteString(fmt.Sprintf("  ▶ %s\n", opt))
		} else {
			sb.WriteString(fmt.Sprintf("    %s\n", opt))
		}
	}
	return sb.String()
}

func (e *Executor) executeSelect(step Step) error {
	prompt := e.parser.Parse(step.Prompt)

	m := selectModel{
		prompt:  prompt,
		options: step.Options,
	}
	p := tea.NewProgram(m)

	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("select failed: %w", err)
	}

	final := result.(selectModel)
	if !final.done {
		return fmt.Errorf("selection cancelled")
	}

	fmt.Printf("  ✔ %s\n\n", final.selected)
	e.parser.Set(step.Variable, final.selected)
	return nil
}
