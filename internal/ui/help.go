package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) helpView() string {
	width := 50
	height := 12

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorInfo).
		Background(lipgloss.Color("236")). // Dark gray bg for modal
		Padding(1, 2)

	var s strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		MarginBottom(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(ColorInfo).
		Bold(true).
		Width(10)

	descStyle := lipgloss.NewStyle().
		Foreground(ColorText)

	row := func(key, desc string) string {
		return fmt.Sprintf("%s%s\n", keyStyle.Render(key), descStyle.Render(desc))
	}

	s.WriteString(titleStyle.Render("Help & Controls"))
	s.WriteString("\n\n")

	s.WriteString(row("q / Esc", "Quit application"))
	s.WriteString(row("r", "Refresh dashboard"))
	s.WriteString(row("?", "Close help"))

	s.WriteString("\n")
	s.WriteString(StyleDim.Render("GitDash v0.1.0"))

	// Center the modal
	// We need to calculate margins to center it in m.Width x m.Height
	// Lipgloss Place can do this effectively if we render the modal as a string first
	modal := style.Render(s.String())

	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, modal)
}
