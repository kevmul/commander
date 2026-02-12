package commandwizard

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "github.com/kevmul/cmdr/internal/workflow"
	"github.com/kevmul/cmdr/internal/styles"
)

const (
	Input = iota
	Select
	Command
)

type Model struct {
	cursor int
	Type   int

	prompt textinput.Model
}

func New() Model {
	return Model{
		cursor: 0,
		prompt: textinput.New(),
	}
}

func NewUpdate() Model {
	return Model{
		cursor: 0,
		prompt: textinput.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			}
		case "enter":
			m.Type = m.cursor
		}
	}
	return m, nil
}

func (m Model) View() string {
	title := styles.TitleStyle.Render("Command Type")
	subtitle := styles.SubtitleStyle.Render("Choose the type of command you want to create")
	menu := m.renderSelectType()

	var output string

	switch m.Type {
	case Input:
		output = styles.InfoStyle.Render("Choose a Prompt Template to use for this command.\n")
		m.prompt.Focus()
		output += m.prompt.View()
	case Select:
		output += " (Select)"
	case Command:
		output += " (Command)"
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		menu,
		output,
	)
}

func (m Model) renderSelectType() string {
	options := []string{"Input", "Select", "Command"}
	var b strings.Builder
	for i, option := range options {
		if m.cursor == i {
			b.WriteString(styles.SelectedItemStyle.Render("> "+option) + "\n")
		} else {
			b.WriteString(styles.NormalItemStyle.Render("  "+option) + "\n")
		}
	}
	return b.String()
}
