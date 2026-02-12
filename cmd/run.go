package cmd

import (
	"fmt"
	"github.com/kevmul/cmdr/internal/workflow"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [workflow-name]",
	Short: "Run a workflow by name",
	Long:  "Execute a workflow directly by name, or show selection menu if no name provided",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := workflow.NewStore()
		if err != nil {
			return err
		}

		// If no workflow name provided, show list (same as cmdr list)
		if len(args) == 0 {
			return listCmd.RunE(cmd, args)
		}

		workflowName := args[0]

		// Check if workflow exists
		if !store.Exists(workflowName) {
			return fmt.Errorf("workflow '%s' not found", workflowName)
		}

		// Load the workflow
		wf, err := store.Load(workflowName)
		if err != nil {
			return err
		}

		// Execute the workflow
		executor := workflow.NewExecutor()
		return executor.Execute(wf)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

