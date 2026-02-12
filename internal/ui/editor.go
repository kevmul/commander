package ui

import (
	"fmt"
	"strings"
	"github.com/kevmul/cmdr/internal/workflow"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type editorModel struct {
	workflow    *workflow.Workflow
	nameInput   textinput.Model
	descInput   textinput.Model
	focusIndex  int
	isEdit      bool
	steps       []workflow.Step
	currentStep int
	stepMode    bool // true when editing a step
	saved       bool
}

func NewEditorModel(wf *workflow.Workflow, isEdit bool) tea.Model {
	m := &editorModel{
		isEdit:   isEdit,
		steps:    []workflow.Step{},
		stepMode: false,
	}

	if isEdit && wf != nil {
		m.workflow = wf
		m.steps = wf.Steps
	} else {
		m.workflow = &workflow.Workflow{}
	}

	m.nameInput = textinput.New()
	m.nameInput.Placeholder = "workflow-name"
	m.nameInput.Focus()
	m.nameInput.CharLimit = 50
	m.nameInput.Width = 50
	if m.workflow.Name != "" {
		m.nameInput.SetValue(m.workflow.Name)
	}

	m.descInput = textinput.New()
	m.descInput.Placeholder = "Description (optional)"
	m.descInput.CharLimit = 100
	m.descInput.Width = 50
	if m.workflow.Description != "" {
		m.descInput.SetValue(m.workflow.Description)
	}

	return m
}

func (m *editorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *editorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+s":
			m.workflow.Name = m.nameInput.Value()
			m.workflow.Description = m.descInput.Value()
			m.workflow.Steps = m.steps
			m.saved = true
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == 2 {
				// TODO: Add step editor
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > 2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 2
			}

			cmds := make([]tea.Cmd, 2)
			for i := 0; i <= 1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.getInput(i).Focus()
				} else {
					m.getInput(i).Blur()
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *editorModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.descInput, cmds[1] = m.descInput.Update(msg)

	return tea.Batch(cmds...)
}

func (m *editorModel) getInput(index int) *textinput.Model {
	switch index {
	case 0:
		return &m.nameInput
	case 1:
		return &m.descInput
	}
	return nil
}

func (m *editorModel) View() string {
	var b strings.Builder

	title := "Create Workflow"
	if m.isEdit {
		title = "Edit Workflow"
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	b.WriteString(m.inputView(0, "Name:", m.nameInput))
	b.WriteString(m.inputView(1, "Description:", m.descInput))

	b.WriteString("\n")
	stepsLabel := "Steps:"
	if m.focusIndex == 2 {
		stepsLabel = focusedStyle.Render("> " + stepsLabel)
	} else {
		stepsLabel = blurredStyle.Render("  " + stepsLabel)
	}
	b.WriteString(stepsLabel + "\n")

	if len(m.steps) == 0 {
		b.WriteString(blurredStyle.Render("  (no steps added yet)\n"))
	} else {
		for i, step := range m.steps {
			stepDesc := fmt.Sprintf("  %d. [%s] %s", i+1, step.Type, step.Prompt)
			if step.Exec != "" {
				stepDesc = fmt.Sprintf("  %d. [%s] %s", i+1, step.Type, step.Exec)
			}
			b.WriteString(blurredStyle.Render(stepDesc) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press enter on Steps to add/edit • ctrl+s to save • esc to cancel"))

	return b.String()
}

func (m *editorModel) inputView(index int, label string, input textinput.Model) string {
	var b strings.Builder

	if m.focusIndex == index {
		b.WriteString(focusedStyle.Render("> " + label))
	} else {
		b.WriteString(blurredStyle.Render("  " + label))
	}
	b.WriteString("\n  ")
	b.WriteString(input.View())
	b.WriteString("\n\n")

	return b.String()
}

func (m *editorModel) GetWorkflow() (*workflow.Workflow, bool) {
	if m.saved {
		return m.workflow, true
	}
	return nil, false
}

func RunEditorUI(wf *workflow.Workflow, isEdit bool) (*workflow.Workflow, error) {
	m := NewEditorModel(wf, isEdit)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	editorModel := finalModel.(*editorModel)
	workflow, saved := editorModel.GetWorkflow()
	if !saved {
		return nil, fmt.Errorf("workflow creation cancelled")
	}

	return workflow, nil
}
