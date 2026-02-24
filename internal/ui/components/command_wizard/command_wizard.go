package commandwizard

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	// "github.com/kevmul/cmdr/internal/workflow"
	"github.com/kevmul/cmdr/internal/styles"
	"github.com/kevmul/cmdr/internal/workflow"
)

const (
	Input = iota
	Select
	Command
)

const (
	TypeSelection = iota
	Fields
	AdvancedOptions
)

type Model struct {
	cursor int

	// Command type: Input, Select, Command
	Type int

	// Command steps
	steps []workflow.Step

	// Input filed
	promptInput  textinput.Model
	varInput     textinput.Model
	commandInput textinput.Model

	Step int
}

func New() Model {
	pi := textinput.New()
	pi.Placeholder = "Enter prompt..."
	pi.Focus()

	ci := textinput.New()
	ci.Placeholder = "Enter command..."
	ci.Width = 50
	ci.CharLimit = 100

	return Model{
		cursor:       0,
		Step:         TypeSelection,
		promptInput:  pi,
		commandInput: ci,
	}
}

func NewUpdate() Model {
	return Model{
		cursor:      0,
		promptInput: textinput.New(),
		Step:        TypeSelection,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			}

		case "shift+tab":
			if m.Step > TypeSelection {
				m.Step--
			}
		case "tab":
			if m.Step < AdvancedOptions {
				m.Step++
			}
		case "enter":
			switch m.Step {
			case TypeSelection:
				m.Type = m.cursor
				m.Step = Fields
			case Fields:
				m.Step = AdvancedOptions
			case AdvancedOptions:
				// Finalize command creation
			}
		}
	}

	switch m.Step {
	case TypeSelection:
		// No input to update
	case Fields:
		switch m.Type {
		case Input:
			m.promptInput, _ = m.promptInput.Update(msg)
			// m.promptText = m.promptInput.Value()
		case Command:
			m.commandInput, _ = m.commandInput.Update(msg)
			// m.commandText = m.commandInput.Value()
		}
	}
	return m, nil
}

func (m Model) View() string {

	switch m.Step {
	case TypeSelection:
		return m.renderTypeOptions()
	case Fields:
		return m.renderFields()
	case AdvancedOptions:
		return m.renderAdvancedOptions()
	default:
		return "Unknown step"
	}
}

/* ==================================================
 * Type Selection
 * =================================================*/

func (m Model) renderTypeOptions() string {

	title := styles.TitleStyle.Render("Command Type")
	subtitle := styles.SubtitleStyle.Render("Choose the type of command you want to create")
	menu := m.renderSelectType()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		menu,
	)
}

func (m Model) renderSelectType() string {
	options := []string{"Input", "Select", "Command"}
	var b strings.Builder
	for i, option := range options {
		if m.cursor == i {
			b.WriteString(styles.SelectedItemStyle.Render("> "+option) + "\n")
		} else {
			b.WriteString(styles.NormalItemStyle.Render("  "+option) + "\n")
		}
	}
	return b.String()
}

/* ==================================================
 * Fields
 * =================================================*/

func (m Model) renderFields() string {
	title := styles.TitleStyle.Render("Setup: ")
	subtitle := styles.SubtitleStyle.Render("Define the fields for your command")

	output := ""
	switch m.Type {
	case Input:
		// Render input fields
		title += "Input Command"
		output = m.renderInputFields()
	case Select:
		// Render select fields
		title += "Select Command"
		output = m.renderSelectFields()
	case Command:
		// Render command fields
		title += "Command String"
		output = m.renderCommandFields()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		output,
	)
}

func (m Model) renderInputFields() string {
	// Render input fields
	return m.promptInput.View()
}

func (m Model) renderSelectFields() string {
	// Render select fields
	return "Select Fields"
}

func (m Model) renderCommandFields() string {
	// Render command fields
	return "Command Fields"
}

/* ==================================================
 * Advanced Options
 * =================================================*/

func (m Model) renderAdvancedOptions() string {
	title := styles.TitleStyle.Render("Advanced Options")
	subtitle := styles.SubtitleStyle.Render("Configure advanced options for your command")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
	)
}
