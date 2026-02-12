package workflowform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevmul/cmdr/internal/messages"
	"github.com/kevmul/cmdr/internal/styles"
)

func (m Model) renderStepCommands() string {

	title := styles.TitleStyle.Render("Commands")
	subtitle := styles.SubtitleStyle.Width(styles.ModalWidth - 4).Render("Add the actual command(s) that will be executed. You can add multiple commands if needed, and they will be run sequentially. For example, you might want to run a build command followed by a test command.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"Press [n] to add a new command.",
	)
}

func (m Model) updateStepCommands(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// Add a new command to the list
			return m, func() tea.Msg {
				return messages.ShowModalMsg{
					ModalType: "create_command",
				}
			}
		}
	}
	return m, nil
}
