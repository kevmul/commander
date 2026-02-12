package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Color
	Primary      = lipgloss.Color("#9ECE6A")
	Secondary    = lipgloss.Color("#BB9AF7")
	Tertiary     = lipgloss.Color("#7AA2F7")
	Success      = lipgloss.Color("#10B981")
	Error        = lipgloss.Color("#EF4444")
	Warning      = lipgloss.Color("#F59E0B")
	Muted        = lipgloss.Color("#6B7280")
	Text         = lipgloss.Color("#D4D4D4")
	Background   = lipgloss.Color("#1E1E1E")
	CustomBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
		Right:       "│",
	}

	// Sizes
	PaddingSmall  = 1
	PaddingMedium = 2
	PaddingLarge  = 3

	MarginSmall  = 1
	MarginMedium = 2
	MarginLarge  = 3

	ModalWidth  = 64
	ModalHeight = 12

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Text).
			Italic(true).
			MarginBottom(1)

	MutedTextStyle = lipgloss.NewStyle().
			Foreground(Muted)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Secondary)

		// Button Styles
	ButtonStyle = lipgloss.NewStyle().
			Background(Muted).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2).
			Bold(true)

	ActiveButtonStyle = lipgloss.NewStyle().
				Background(Secondary).
				Foreground(lipgloss.Color("#FFFFFF")).
				Underline(true).
				Padding(0, 2).
				Bold(true)

	// Input styles
	FocusedInputStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(0, 1)

	BlurredInputStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Muted).
				Padding(0, 1)

	// List styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Secondary).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle()

	// Command Styles
	CommandStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1).
			MarginBottom(1)

	// Modal Styles
	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderLeft(true).
			BorderRight(true).
			BorderTop(false).
			BorderBottom(true).
			BorderForeground(Secondary).
			Padding(0, 1).
			Width(ModalWidth)

	ModalWithScrollStyle = lipgloss.NewStyle().
				BorderStyle(CustomBorder).
				BorderLeft(true).
				BorderRight(false).
				BorderTop(false).
				BorderBottom(true).
				BorderForeground(Secondary).
				PaddingLeft(1).
				Width(ModalWidth)

	ModalTitleStyle = lipgloss.NewStyle().
			Foreground(Secondary)
)
