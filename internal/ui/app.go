package ui

import (
	"github.com/kevmul/cmdr/internal/messages"
	"github.com/kevmul/cmdr/internal/styles"
	"github.com/kevmul/cmdr/internal/ui/components/modal"
	"github.com/kevmul/cmdr/internal/utils"
	"github.com/kevmul/cmdr/internal/workflow"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 0)
)

type keyMap struct {
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Run    key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "run"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type workflowItem struct {
	workflow.Workflow
}

func (i workflowItem) Title() string       { return i.Name }
func (i workflowItem) Description() string { return i.Workflow.Description }
func (i workflowItem) FilterValue() string { return i.Name }

type mainModel struct {
	list     list.Model
	store    *workflow.Store
	selected *workflow.Workflow
	action   string // "run", "edit", "delete", "create", ""

	showModal bool
	modal     *modal.Model

	width  int
	height int

	ready bool
}

func NewMainModel(store *workflow.Store) (tea.Model, error) {

	workflows, err := store.List()
	if err != nil {
		return nil, err
	}

	items := make([]list.Item, len(workflows))
	for i, w := range workflows {
		items[i] = workflowItem{w}
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(styles.Primary).BorderLeftForeground(styles.Primary)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(styles.Tertiary).BorderLeftForeground(styles.Primary)

	l := list.New(items, d, 0, 0)
	l.Title = "Command Runner"
	l.SetShowTitle(false)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowHelp(false)

	return &mainModel{
		list:  l,
		store: store,
	}, nil
}

func (m *mainModel) Init() tea.Cmd {
	return nil
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4)
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		if m.showModal {
			// Let the modal handle key messages
			break
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.New):
			m.action = "create"
			m.showModal = true
			m.modal = modal.NewCreateWorkflow(m.store)

			return m, nil

		case key.Matches(msg, keys.Edit):
			if item, ok := m.list.SelectedItem().(workflowItem); ok {
				m.selected = &item.Workflow
				m.action = "edit"
				return m, tea.Quit
			}

		case key.Matches(msg, keys.Delete):
			if item, ok := m.list.SelectedItem().(workflowItem); ok {
				m.selected = &item.Workflow
				m.action = "delete"
				return m, tea.Quit
			}

		case key.Matches(msg, keys.Run):
			if item, ok := m.list.SelectedItem().(workflowItem); ok {
				m.selected = &item.Workflow
				m.action = "run"
				return m, tea.Quit
			}
		}

	case messages.ModalClosedMsg:
		m.showModal = false
		m.modal = nil

	case messages.ShowModalMsg:
		m.showModal = true
		switch msg.ModalType {
		case "create_command":
			m.modal = modal.NewCreateCommand()
			// case "edit_command":
			// 	if item, ok := m.list.SelectedItem().(workflowItem); ok {
			// 		m.selected = &item.Workflow
			// 		m.modal = modal.NewEditWorkflow(m.store, m.selected)
			// 	}
		}

	case messages.WorkflowCreateMsg:
		newWorkflow := workflow.Workflow{
			Name:        msg.Name,
			Description: msg.Description,
		}
		err := m.store.Save(&newWorkflow)
		if err != nil {
			// Handle error (e.g., show an error modal)
			break
		}
		m.showModal = false
		m.modal = nil

		// Refresh the workflow list
		workflows, err := m.store.List()
		if err != nil {
			break
		}
		items := make([]list.Item, len(workflows))
		for i, w := range workflows {
			items[i] = workflowItem{w}
		}
		m.list.SetItems(items)

	}
	if m.showModal && m.modal != nil {
		*m.modal, cmd = m.modal.Update(msg)
		cmds = append(cmds, cmd)
	}

	if !m.showModal {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *mainModel) View() string {

	helpText := helpStyle.Render("[n] New  [e] Edit  [d] Delete  [↵] Run  [q] Quit")
	currentView := m.list.View() + "\n" + helpText

	if m.showModal && m.modal != nil {
		currentView = utils.RenderWithModal(m.height, m.width, currentView, m.modal.View())
	}

	return currentView
}

func (m *mainModel) GetAction() (string, *workflow.Workflow) {
	return m.action, m.selected
}

func RunMainUI(store *workflow.Store) error {
	m, err := NewMainModel(store)
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}
