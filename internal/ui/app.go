package ui

import (
	"fmt"

	"github.com/kevmul/cmdr/internal/styles"
	"github.com/kevmul/cmdr/internal/workflow"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Padding(1, 0)
)

type keyMap struct {
	Run  key.Binding
	Quit key.Binding
}

var keys = keyMap{
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

func (i workflowItem) Title() string {
	key := styles.MutedTextStyle.Render(fmt.Sprintf("(%s)", i.Key))
	return fmt.Sprintf("%s %s", i.Name, key)
}
func (i workflowItem) Description() string { return i.Workflow.Description }
func (i workflowItem) FilterValue() string { return i.Name }

type mainModel struct {
	list     list.Model
	store    *workflow.Store
	selected *workflow.Workflow
	action   string // "run", "edit", "delete", "create", ""

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
		ready: false,
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
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.Run):
			if item, ok := m.list.SelectedItem().(workflowItem); ok {
				m.selected = &item.Workflow
				m.action = "run"
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m *mainModel) View() string {

	helpText := helpStyle.Render("[n] New  [e] Edit  [d] Delete  [↵] Run  [q] Quit")
	currentView := m.list.View() + "\n" + helpText

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
