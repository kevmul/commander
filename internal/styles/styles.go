package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Color
	Primary    = lipgloss.Color("#9ECE6A")
	Secondary  = lipgloss.Color("#BB9AF7")
	Tertiary   = lipgloss.Color("#7AA2F7")
	Success    = lipgloss.Color("#10B981")
	Error      = lipgloss.Color("#EF4444")
	Warning    = lipgloss.Color("#F59E0B")
	Muted      = lipgloss.Color("#6B7280")
	Text       = lipgloss.Color("#D4D4D4")
	Background = lipgloss.Color("#1E1E1E")
	Clear      = lipgloss.Color("")

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

	InputStyle = lipgloss.NewStyle()

	HelpTextStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	CursorStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	// List styles
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Secondary).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle()

	ListItemStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(false).
			MarginLeft(1).
			BorderForeground(Clear).
			Foreground(Text).
			Padding(0, 2).
			MarginBottom(1)

	ListItemSelectedStyle = ListItemStyle.
				BorderLeft(true).
				BorderForeground(Secondary).
				Foreground(Primary).
				Bold(true).
				Padding(0, 1)

	ListItemTitleStyle = lipgloss.NewStyle().
				Foreground(Text)
	ListItemDescriptionStyle = lipgloss.NewStyle().
					Foreground(Muted)
	ListItemSelectedTitleStyle = lipgloss.NewStyle().
					Foreground(Primary).
					Bold(true)
	ListItemSelectedDescriptionStyle = lipgloss.NewStyle().
						Foreground(Secondary).
						Bold(true)

	// Command Styles
	CommandStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1).
			MarginBottom(1)
)
