package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/kevmul/cmdr/internal/workflow"

	"github.com/spf13/cobra"
)

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

		// Build options for gum choose
		var options []string
		workflowMap := make(map[string]*workflow.Workflow)
		
		for i := range workflows {
			wf := &workflows[i]
			option := wf.Name
			if wf.Description != "" {
				option = fmt.Sprintf("%s - %s", wf.Name, wf.Description)
			}
			options = append(options, option)
			workflowMap[wf.Name] = wf
		}

		// Use gum to select
		gumArgs := append([]string{"choose"}, options...)
		gumCmd := exec.Command("gum", gumArgs...)
		gumCmd.Stderr = os.Stderr

		output, err := gumCmd.Output()
		if err != nil {
			return fmt.Errorf("selection cancelled or gum not installed")
		}

		selected := strings.TrimSpace(string(output))
		
		// Extract workflow name from selection
		workflowName := strings.Split(selected, " - ")[0]
		
		selectedWorkflow := workflowMap[workflowName]
		if selectedWorkflow == nil {
			return fmt.Errorf("workflow not found")
		}

		// Execute the selected workflow
		executor := workflow.NewExecutor()
		return executor.Execute(selectedWorkflow)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

