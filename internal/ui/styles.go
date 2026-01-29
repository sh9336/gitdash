package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#7D56F4")
	ColorSecondary = lipgloss.Color("#FAFAFA")
	ColorText      = lipgloss.Color("#E0E0E0")
	ColorSubText   = lipgloss.Color("#9E9E9E")

	ColorSuccess = lipgloss.Color("#04B575") // Green
	ColorWarning = lipgloss.Color("#FFD93D") // Yellow
	ColorError   = lipgloss.Color("#FF6B6B") // Red
	ColorInfo    = lipgloss.Color("#3B8ED0") // Blue

	// Styles
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Background(ColorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	StylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 0).
			MarginRight(0)

	StyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")) // Cyan-ish

	StyleSelected = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	StyleNormal = lipgloss.NewStyle().
			Foreground(ColorText)

	StyleDim = lipgloss.NewStyle().
			Foreground(ColorSubText)
)
