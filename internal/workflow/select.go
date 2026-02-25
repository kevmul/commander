package workflow

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
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

type keyMap struct {
	Up   key.Binding
	Down key.Binding
	Run  key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k, up", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j, down", "Move down"),
	),
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "run"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q, ctrl+c", "quit"),
	),
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, keys.Down):
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case key.Matches(msg, keys.Run):
			m.selected = m.options[m.cursor]
			m.done = true
			return m, tea.Quit
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:\n", m.prompt))
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

	fmt.Printf("\n  ✔  %s\n\n", final.selected)
	e.parser.Set(step.Variable, final.selected)
	return nil
}
