package workflow

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevmul/cmdr/internal/styles"
)

type messageModel struct {
	output  string
	variant string
	done    bool
}

func (m messageModel) Init() tea.Cmd { return nil }

func (m messageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyEsc, tea.KeyCtrlC:
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m messageModel) View() string {
	var color lipgloss.Color
	switch m.variant {
	case "error":
		color = styles.Error
	case "success":
		color = styles.Success
	case "warning":
		color = styles.Warning
	default:
		color = styles.Primary
	}

	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.BlockBorder()).
		Foreground(color).
		Render(m.output)
}

func (e *Executor) executeMessage(step Step) error {
	text := e.parser.Parse(step.Prompt)
	variant := e.parser.Parse(step.Variant)
	m := messageModel{variant: variant, output: text}
	fmt.Println(m.View())
	return nil
}
