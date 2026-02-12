package commandform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevmul/cmdr/internal/styles"
)

func (m Model) renderStepCommands() string {
	title := styles.TitleStyle.Render("Commands")
	subtitle := styles.SubtitleStyle.Width(styles.ModalWidth - 4).Render("Add the actual command(s) that will be executed. You can add multiple commands if needed, and they will be run sequentially. For example, you might want to run a build command followed by a test command.")

	// For testing, create a few dummy commands
	m.commands = []string{
		"npm install",
		"npm run build",
		"npm test",
	}

	// Create a list of command views
	var commandViews []string
	for i, cmd := range m.commands {
		cmdView := styles.CommandStyle.Render(cmd)
		commandViews = append(commandViews, styles.CommandStyle.Render(string(i+1)+".")+cmdView)
	}

	commandsList := lipgloss.JoinVertical(lipgloss.Left, commandViews...)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		commandsList,
	)
}

func (m Model) updateStepCommands(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}
