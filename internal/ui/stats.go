package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sh9336/gitdash/internal/stats"
)

type StatsModel struct {
	Stats *stats.ProjectStats
}

func NewStatsModel(s *stats.ProjectStats) StatsModel {
	return StatsModel{
		Stats: s,
	}
}

func (m StatsModel) View(width int) string {
	var s strings.Builder

	// Header
	s.WriteString(StyleHeader.Render("Project Stats"))
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", width)))
	s.WriteString("\n")

	if m.Stats == nil {
		s.WriteString(StyleDim.Render("   No stats available"))
		return StylePanel.Copy().Width(width).Render(s.String())
	}

	// Total files
	s.WriteString(fmt.Sprintf(" Total Files: %d\n\n", m.Stats.TotalFiles))

	// Languages
	for i, l := range m.Stats.Languages {
		if i >= 5 {
			break
		}

		barWidth := int(l.Percentage / 2) // scale down
		bar := strings.Repeat("█", barWidth)

		line := fmt.Sprintf(" %-10s %5.1f%% %s", l.Name, l.Percentage, bar)
		s.WriteString(line + "\n")
	}

	return StylePanel.Copy().Width(width).Render(s.String())
}
