package modal

import (
	"bytes"
	"strings"

	"github.com/kevmul/cmdr/internal/messages"
	"github.com/kevmul/cmdr/internal/styles"
	"github.com/kevmul/cmdr/internal/ui/components/command_wizard"
	"github.com/kevmul/cmdr/internal/ui/components/confirmation"
	"github.com/kevmul/cmdr/internal/ui/components/workflow_form"
	"github.com/kevmul/cmdr/internal/utils"
	"github.com/kevmul/cmdr/internal/workflow"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModalType int

const (
	// Workflow Modals
	CreateWorkflow ModalType = iota
	UpdateWorkflow
	// Command Modals
	CreateCommand
	UpdateCommand
	// Delete Confirmation Modal
	DeleteConfirmation
	HelpModal
)

type Model struct {
	modalType      ModalType
	createWorkflow *workflowform.Model
	updateWorkflow *workflowform.Model

	createCommand *commandwizard.Model
	updateCommand *workflowform.Model

	deleteConfirmation *confirmation.Model
	// UI
	viewport viewport.Model
	title    string
}

// =========================================
// Workflow Modals
// =========================================

func NewCreateWorkflow(store *workflow.Store) *Model {
	createWorkflow := workflowform.New()
	viewport := viewport.New(0, 4)
	viewport.SetContent(createWorkflow.View())

	return &Model{
		modalType:      CreateWorkflow,
		createWorkflow: &createWorkflow,
		viewport:       viewport,
		title:          "Create New Workflow",
	}
}

func NewUpdateWorkflow(workflowId string) *Model {
	updateWorkflow := workflowform.NewUpdate(workflowId)
	viewport := viewport.New(0, 4)
	viewport.SetContent(updateWorkflow.View())

	return &Model{
		modalType:      UpdateWorkflow,
		updateWorkflow: &updateWorkflow,
		viewport:       viewport,
		title:          "Update Workflow",
	}
}

// =========================================
// Command Modals
// =========================================

func NewCreateCommand() *Model {
	createCommand := commandwizard.New()
	viewport := viewport.New(0, 4)
	viewport.SetContent(createCommand.View())

	return &Model{
		modalType:     CreateCommand,
		createCommand: &createCommand,
		viewport:      viewport,
		title:         "Create New Command",
	}
}

// func NewUpdateCommand(commandId string) *Model {
// 	updateCommand := commandwizard.NewUpdate()
// 	viewport := viewport.New(0, 4)
// 	viewport.SetContent(updateCommand.View())
//
// 	return &Model{
// 		modalType:      UpdateWorkflow,
// 		updateCommand : &updateCommand,
// 		viewport:       viewport,
// 		title:          "Update Command",
// 	}
// }

// =========================================
// Delete Confirmation Modal
// =========================================

func NewDeleteConfirmation(itemId, itemType string) *Model {
	deleteConfirmation := confirmation.New(itemId, itemType)
	viewport := viewport.New(0, 4)
	viewport.SetContent(deleteConfirmation.View())

	return &Model{
		modalType:          DeleteConfirmation,
		deleteConfirmation: &deleteConfirmation,
		viewport:           viewport,
		title:              "Confirm Deletion",
	}
}

func (m Model) Init() tea.Cmd {
	switch m.modalType {
	case DeleteConfirmation:
		return m.deleteConfirmation.Init()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// We might move this to the modal themselves later...
		if msg.String() == "esc" || msg.String() == "q" || msg.String() == "ctrl+c" {
			// Handle closing in the modal itself if needed (e.g. to reset form state), then send message to parent to close the modal
			var cmd tea.Cmd
			// Send a message to parent to close the modal
			return m, tea.Batch(cmd, func() tea.Msg {
				return messages.ModalClosedMsg{}
			})
		}
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.modalType {
	case DeleteConfirmation:
		*m.deleteConfirmation, cmd = m.deleteConfirmation.Update(msg)
		if getLines(m.deleteConfirmation.View()) != m.viewport.TotalLineCount() {
			// Update viewport height if content height changes
			m.viewport.Height = getLines(m.deleteConfirmation.View())
		}
	case CreateWorkflow:
		*m.createWorkflow, cmd = m.createWorkflow.Update(msg)
		if getLines(m.createWorkflow.View()) != m.viewport.TotalLineCount() {
			// Update viewport height if content height changes
			m.viewport.Height = min(styles.ModalHeight, getLines(m.createWorkflow.View())+1)
		}
	case UpdateWorkflow:
		*m.updateWorkflow, cmd = m.updateWorkflow.Update(msg)

	case CreateCommand:
		*m.createCommand, cmd = m.createCommand.Update(msg)
		if getLines(m.createCommand.View()) != m.viewport.TotalLineCount() {
			// Update viewport height if content height changes
			m.viewport.Height = min(styles.ModalHeight, getLines(m.createCommand.View())+1)
		}
	}
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	if _, ok := msg.(tea.KeyMsg); ok {
		// Update viewport content on key events
		m.viewport.SetContent(m.RenderContent())
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	if m.viewport.TotalLineCount() <= m.viewport.Height {
		// No scrollbar needed
		// return styles.ModalStyle.Width(styles.ModalWidth).Render(viewport)
		return lipgloss.JoinVertical(
			lipgloss.Top,
			createBorderTitle(m.title, styles.ModalWidth, false),
			styles.ModalStyle.Render(m.viewport.View()),
		)
	}

	viewport := lipgloss.JoinVertical(
		lipgloss.Top,
		createBorderTitle(m.title, styles.ModalWidth, true),
		styles.ModalWithScrollStyle.Render(m.viewport.View()),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		// styles.ModalWithScrollStyle.Width(styles.ModalWidth).Render(viewport),

		viewport,
		utils.RenderScrollbarForModal(m.viewport),
	)

}

func createBorderTitle(title string, modalWidth int, withScroll bool) string {
	borderChar := styles.CustomBorder.Top
	titleLength := lipgloss.Width(title)
	if titleLength >= modalWidth-2 {
		// Title is too long to fit, return it as is (it will be truncated by the modal)
		return title
	}

	leftBorderLength := 2                                                //
	rightBorderLength := modalWidth - titleLength - leftBorderLength - 2 // 2 for the spaces around the title

	s := styles.CustomBorder.TopLeft +
		strings.Repeat(string(borderChar), leftBorderLength) +
		" " + title + " " +
		strings.Repeat(string(borderChar), rightBorderLength)

	if !withScroll {
		s += styles.CustomBorder.TopRight
	}

	return styles.ModalTitleStyle.Render(s)

}

func (m Model) RenderContent() string {
	switch m.modalType {
	case CreateWorkflow:
		return m.createWorkflow.View()
	case CreateCommand:
		return m.createCommand.View()
	case DeleteConfirmation:
		return m.deleteConfirmation.View()
	}
	return "MODAL"
}

func getLines(s string) int {
	return bytes.Count([]byte(s), []byte{'\n'})
}
