package ui

import (
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

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Command Runner"
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v-4)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.New):
			m.action = "create"
			return m, tea.Quit

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
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *mainModel) View() string {
	helpText := helpStyle.Render("[n] New  [e] Edit  [d] Delete  [↵] Run  [q] Quit")
	return m.list.View() + "\n" + helpText
}

func (m *mainModel) GetAction() (string, *workflow.Workflow) {
	return m.action, m.selected
}

func RunMainUI(store *workflow.Store) (string, *workflow.Workflow, error) {
	m, err := NewMainModel(store)
	if err != nil {
		return "", nil, err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", nil, err
	}

	mainModel := finalModel.(*mainModel)

	action, workflow := mainModel.GetAction()

	return action, workflow, nil
}


