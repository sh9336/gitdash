package ui

import (
	"fmt"
	"strings"

	"github.com/gitdash/gitdash/internal/stats"
)

type StatsModel struct {
	Stats *stats.ProjectStats
}

func NewStatsModel(s *stats.ProjectStats) StatsModel {
	return StatsModel{
		Stats: s,
	}
}

func (m StatsModel) View() string {
	var s strings.Builder

	s.WriteString(StyleHeader.Render("Project Stats"))
	s.WriteString("\n")

	if m.Stats == nil {
		s.WriteString(StyleDim.Render("  No stats available"))
		return StylePanel.Render(s.String())
	}

	// Total files
	s.WriteString(fmt.Sprintf("Total Files: %d\n\n", m.Stats.TotalFiles))

	// Languages
	for i, l := range m.Stats.Languages {
		if i >= 5 {
			break
		}

		barWidth := int(l.Percentage / 2) // scale down
		bar := strings.Repeat("█", barWidth)

		// Color would be nicer if we had a map of lang -> color
		// For now just use distinct colors from lipgloss

		line := fmt.Sprintf("%-10s %5.1f%% %s", l.Name, l.Percentage, bar)
		s.WriteString(line + "\n")
	}

	return StylePanel.Render(s.String())
}
