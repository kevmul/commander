package workflowform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevmul/cmdr/internal/styles"
)

func (m Model) renderStepDescription() string {
	title := styles.TitleStyle.Render("Command Description")
	subtitle := styles.SubtitleStyle.Render("Describe what this command does. This is optional but can be helpful for users to understand the purpose of the command.")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		m.descriptionInput.View(),
	)
}

func (m Model) updateStepDescription(msg tea.Msg) (Model, tea.Cmd) {
	m.descriptionInput, _ = m.descriptionInput.Update(msg)
	return m, nil
}
