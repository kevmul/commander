package commandform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevmul/cmdr/internal/styles"
)

func (m Model) renderStepName() string {
	title := styles.TitleStyle.Render("Command Name")
	subtitle := styles.SubtitleStyle.Render("This name should contain no spaces and be unique within the workflow.")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		m.nameInput.View(),
	)
}

func (m Model) updateStepName(msg tea.Msg) (Model, tea.Cmd) {
	m.nameInput, _ = m.nameInput.Update(msg)
	return m, nil
}
