package cmd

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevmul/cmdr/internal/ui"
	"github.com/kevmul/cmdr/internal/workflow"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cmdr",
	Short: "A command runner for multi-step dev workflows",
	Long: `cmdr is a CLI tool for creating and running multi-step interactive workflows.
Store your common dev tasks as reusable workflows with inputs, selections, and commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := workflow.NewStore()
		if err != nil {
			return err
		}

		m, err := ui.NewMainModel(store)
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err = p.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
