package cmd

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevmul/cmdr/internal/workflow"
	"github.com/spf13/cobra"
)

// ─── Inline workflow selector model ───────────────────────────────────────────

type workflowSelectModel struct {
	workflows []workflow.Workflow
	cursor    int
	selected  *workflow.Workflow
	done      bool
}

func (m workflowSelectModel) Init() tea.Cmd { return nil }

func (m workflowSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.workflows)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			m.selected = &m.workflows[m.cursor]
			m.done = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m workflowSelectModel) View() string {
	var sb strings.Builder
	sb.WriteString("? Select a workflow:\n")
	for i, wf := range m.workflows {
		if i == m.cursor {
			if wf.Description != "" {
				sb.WriteString(fmt.Sprintf("  ▶ %s — %s\n", wf.Name, wf.Description))
			} else {
				sb.WriteString(fmt.Sprintf("  ▶ %s\n", wf.Name))
			}
		} else {
			if wf.Description != "" {
				sb.WriteString(fmt.Sprintf("    %s — %s\n", wf.Name, wf.Description))
			} else {
				sb.WriteString(fmt.Sprintf("    %s\n", wf.Name))
			}
		}
	}
	sb.WriteString("\n  ↑/↓ to move · enter to select · esc to cancel\n")
	return sb.String()
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

		m := workflowSelectModel{workflows: workflows}
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

		fmt.Printf("  ✔ %s\n\n", final.selected.Name)

		executor := workflow.NewExecutor()
		return executor.Execute(final.selected)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
