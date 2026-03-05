package cmd

import (
	"fmt"
	// "strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	// "github.com/charmbracelet/lipgloss"
	// "github.com/kevmul/cmdr/internal/styles"
	"github.com/kevmul/cmdr/internal/styles"
	"github.com/kevmul/cmdr/internal/workflow"
	"github.com/spf13/cobra"
)

// ─── Inline workflow selector model ───────────────────────────────────────────

type workflowSelectModel struct {
	workflows []workflow.Workflow
	selected  *workflow.Workflow
	list      list.Model
	listItems []list.Item
	done      bool
}

type keyMap struct {
	Up   key.Binding
	Down key.Binding
	Run  key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k, up", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j, down", "Move down"),
	),
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "run"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q, ctrl+c", "quit"),
	),
}

func (m workflowSelectModel) Init() tea.Cmd { return nil }

type workflowItem struct {
	workflow.Workflow
}

func (i workflowItem) Title() string {
	key := styles.MutedTextStyle.Render(fmt.Sprintf("(%s)", i.Key))
	return fmt.Sprintf("%s %s", i.Name, key)
}
func (i workflowItem) Description() string { return i.Workflow.Description }
func (i workflowItem) FilterValue() string { return i.Name }

func NewWorkflowSelectModel(workflows []workflow.Workflow) workflowSelectModel {

	items := make([]list.Item, len(workflows))
	for i, item := range workflows {
		items[i] = workflowItem{item}
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
	return workflowSelectModel{
		workflows: workflows,
		listItems: items,
		list:      l,
		selected:  nil,
		done:      false,
	}
}

func (m workflowSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		const itemHeight = 2
		const maxVisibleItems = 10
		h := min(len(m.listItems)*itemHeight+4, maxVisibleItems*itemHeight+4)
		m.list.SetSize(msg.Width, h)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Run):
			m.selected = &m.workflows[m.list.Index()]
			m.done = true
			return m, tea.Quit
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}
	}
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m workflowSelectModel) View() string {
	// Returning empty string on done causes BubbleTea to overwrite
	// the list with nothing on its final render — cleanly hiding it.
	if m.done {
		return ""
	}

	return m.list.View()
}

// ─── Command ──────────────────────────────────────────────────────────────────

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List and select a workflow to run",
	Long:  "Display all available workflows in a selection menu and run the chosen one",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := workflow.NewStore()
		if err != nil {
			return err
		}

		workflows, err := store.List()
		if err != nil {
			return err
		}

		if len(workflows) == 0 {
			fmt.Println("No workflows found. Run 'cmdr' to create one!")
			return nil
		}

		m := NewWorkflowSelectModel(workflows)
		p := tea.NewProgram(m)

		result, err := p.Run()
		if err != nil {
			return fmt.Errorf("selection failed: %w", err)
		}

		final := result.(workflowSelectModel)
		if !final.done || final.selected == nil {
			fmt.Println("Cancelled.")
			return nil
		}

		executor := workflow.NewExecutor()
		return executor.Execute(final.selected)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
