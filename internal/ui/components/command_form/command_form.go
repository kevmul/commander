package commandform

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	stepName         = iota // Input for command name
	stepDescription         // Input for command description
	stepArguments           // Input for command arguments
	stepConfirmation        // Confirm and create command
)

type Model struct {
	currentStep      int
	nameInput        textinput.Model
	descriptionInput textinput.Model
	focusIndex       int
	commands         []string
}

func New() Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "command-name"
	nameInput.Focus()
	nameInput.CharLimit = 50
	nameInput.Width = 50

	descInput := textinput.New()
	descInput.Placeholder = "Description (optional)"
	descInput.CharLimit = 100
	descInput.Width = 50
	return Model{
		currentStep:      0,
		nameInput:        nameInput,
		descriptionInput: descInput,
	}
}

func NewUpdate(commandId string) Model {
	nameInput := textinput.New()
	nameInput.Placeholder = "command-name"
	nameInput.Focus()
	nameInput.CharLimit = 50
	nameInput.Width = 50

	descInput := textinput.New()
	descInput.Placeholder = "Description (optional)"
	descInput.CharLimit = 100
	descInput.Width = 50

	return Model{
		currentStep:      0,
		nameInput:        nameInput,
		descriptionInput: descInput,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			switch m.currentStep {
			case stepName:
				m.currentStep = stepDescription
				m.nameInput.Blur()
				m.descriptionInput.Focus()
			case stepDescription:
				m.currentStep = stepArguments
			case stepArguments:
				// m.currentStep = stepConfirmation
			}

		case "tab":
			switch m.currentStep {
			case stepName:
				if m.nameInput.Value() == "" {
					break
				}
				m.currentStep = stepDescription
				m.nameInput.Blur()
				m.descriptionInput.Focus()
			case stepDescription:
				if m.descriptionInput.Value() == "" {
					break
				}
				m.currentStep = stepArguments
			}

		case "shift+tab":
			switch m.currentStep {
			case stepDescription:
				m.currentStep = stepName
				m.descriptionInput.Blur()
				m.nameInput.Focus()
			case stepArguments:
				m.currentStep = stepDescription
				// m.argumentsInput.Blur()
				m.descriptionInput.Focus()
				// case stepConfirmation:
				// 	m.currentStep = stepArguments
				// 	m.confirmationInput.Blur()
				// 	m.argumentsInput.Focus()
			}

		}
	}

	switch m.currentStep {
	case stepName:
		m, cmd = m.updateStepName(msg)
	case stepDescription:
		m, cmd = m.updateStepDescription(msg)
	case stepArguments:
		m, cmd = m.updateStepCommands(msg)
		// case stepConfirmation:
		// 	return m.updateStepConfirmation(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.currentStep {
	case stepName:
		return m.renderStepName()
	case stepDescription:
		return m.renderStepDescription()
	case stepArguments:
		return m.renderStepCommands()
	// case stepConfirmation:
	// 	return m.renderStepConfirmation()
	default:
		return "Unknown step"
	}
}
