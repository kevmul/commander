package workflowform

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevmul/cmdr/internal/styles"
	"github.com/kevmul/cmdr/internal/workflow"
)

func (m Model) renderStepConfirmation() string {
	title := styles.TitleStyle.Render("Confirm")
	subtitle := styles.SubtitleStyle.Width(styles.ModalWidth - 4).Render("Please review the details of your command before confirming. If everything looks good, click 'Confirm' to create your command. If you need to make changes, you can go back to the previous steps.")

	jsonStr := `

{
  "name": "Example2",
  "description": "Example Command",
  "steps": [
    {
      "type": "input",
      "prompt": "Enter your name:",
      "variable": "user_name"
    },
    {
      "type": "select",
      "prompt": "Choose the time of day:",
      "options": ["Morning", "Afternoon", "Evening"],
      "variable": "time_of_day"
    },
    {
      "type": "command",
      "command": "echo Good {{time_of_day}}, {{user_name}}!",
      "capture_output": false
    }
  ]
}`

	// Parse the json

	var workflow workflow.Workflow

	if err := json.Unmarshal([]byte(jsonStr), &workflow); err != nil {
		return fmt.Sprintf("Error parsing JSON: %v", err)
	}

	// Loop through the JSON
	var commandsList strings.Builder
	for i, step := range workflow.Steps {
		// var stepDisplay string
		var stepDisplay strings.Builder
		switch step.Type {
		case "input":
			fmt.Fprintf(&stepDisplay, "%d. [Input] Prompt: %s, Variable: %s", i+1, step.Prompt, step.Variable)
		case "select":
			fmt.Fprintf(&stepDisplay, "%d. [Select] Prompt: %s, Variable: %s", i+1, step.Prompt, step.Variable)
			for _, option := range step.Options {
				fmt.Fprintf(&stepDisplay, "\n   - Option: %s", option)
			}
		case "command":
			fmt.Fprintf(&stepDisplay, "%d. [Command] Command: %s", i+1, step.Command)
		default:
			fmt.Fprintf(&stepDisplay, "%d. [Unknown Step Type]", i+1)
		}

		commandsList.WriteString(lipgloss.NewStyle().Render(stepDisplay.String()) + "\n")
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		commandsList.String(),
	)
}

func (m Model) updateStepConfirmation(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}
