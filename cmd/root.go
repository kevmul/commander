package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/kevmul/cmdr/internal/ui"
	"github.com/kevmul/cmdr/internal/workflow"
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

		action, selected, err := ui.RunMainUI(store)
		if err != nil {
			return err
		}

		switch action {
		case "run":
			if selected != nil {
				executor := workflow.NewExecutor()
				return executor.Execute(selected)
			}
		case "create":
			wf, err := ui.RunEditorUI(nil, false)
			if err != nil {
				fmt.Println("Cancelled")
				return nil
			}
			if err := store.Save(wf); err != nil {
				return err
			}
			fmt.Printf("✅ Workflow '%s' created successfully!\n", wf.Name)
		case "edit":
			if selected != nil {
				wf, err := ui.RunEditorUI(selected, true)
				if err != nil {
					fmt.Println("Cancelled")
					return nil
				}
				if err := store.Save(wf); err != nil {
					return err
				}
				fmt.Printf("✅ Workflow '%s' updated successfully!\n", wf.Name)
			}
		case "delete":
			if selected != nil {
				if err := store.Delete(selected.Name); err != nil {
					return err
				}
				fmt.Printf("✅ Workflow '%s' deleted successfully!\n", selected.Name)
			}
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
