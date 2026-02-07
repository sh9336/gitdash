package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) helpView() string {
	width := 60
	height := 18

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Background(lipgloss.Color("236")).
		Padding(1, 2)

	var s strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).MarginBottom(1)
	keyStyle := lipgloss.NewStyle().Foreground(ColorInfo).Bold(true).Width(12)
	descStyle := lipgloss.NewStyle().Foreground(ColorText)
	guideStyle := lipgloss.NewStyle().Foreground(ColorSubText).Italic(true)

	row := func(key, desc string) string {
		return fmt.Sprintf("%s%s\n", keyStyle.Render(key), descStyle.Render(desc))
	}

	s.WriteString(titleStyle.Render("GitDash - Command Guide"))
	s.WriteString("\n\n")

	s.WriteString(row("Tab", "Toggle focus (General / Branches)"))
	s.WriteString(row("↑/↓ / k/j", "Scroll OR Inspect highlighted branch"))
	s.WriteString(row("f", "Force Checkout (Discard local changes)"))
	s.WriteString(row("r", "Hard Refresh dashboard"))
	s.WriteString(row("?", "Close this menu"))
	s.WriteString(row("q / Esc", "Quit application"))

	s.WriteString("\n")
	s.WriteString(guideStyle.Render("🔍 Inspection Mode is automatic - move the cursor over\nAny branch to view its commits and stats instantly.\nYour working directory is always protected."))

	s.WriteString("\n\n")
	s.WriteString(StyleDim.Render("Built by DeepMind | Press Esc to return"))

	modal := style.Render(s.String())
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, modal)
}
