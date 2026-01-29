package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gitdash/gitdash/internal/git"
)

type WorkDirModel struct {
	Status *git.WorkingDirStatus
}

func NewWorkDirModel(status *git.WorkingDirStatus) WorkDirModel {
	return WorkDirModel{
		Status: status,
	}
}

func (m WorkDirModel) View() string {
	var s strings.Builder

	s.WriteString(StyleHeader.Render("Working Directory"))
	s.WriteString("\n")

	if len(m.Status.Files) == 0 {
		s.WriteString(StyleDim.Render("  Working directory clean"))
		return StylePanel.Render(s.String())
	}

	// Summaries
	// Summaries
	if m.Status.Modified > 0 {
		s.WriteString(fmt.Sprintf("%s Modified: %d\n", lipgloss.NewStyle().Foreground(ColorWarning).Render("●"), m.Status.Modified))
	}
	if m.Status.Staged > 0 {
		s.WriteString(fmt.Sprintf("%s Staged: %d\n", lipgloss.NewStyle().Foreground(ColorSuccess).Render("✓"), m.Status.Staged))
	}
	if m.Status.Untracked > 0 {
		s.WriteString(fmt.Sprintf("%s Untracked: %d\n", lipgloss.NewStyle().Foreground(ColorError).Render("?"), m.Status.Untracked))
	}

	s.WriteString("\n")

	// List files (Top 10 maybe?)
	// Sort by status then name
	sortedFiles := make([]git.FileStatus, len(m.Status.Files))
	copy(sortedFiles, m.Status.Files)

	sort.Slice(sortedFiles, func(i, j int) bool {
		if sortedFiles[i].Status != sortedFiles[j].Status {
			return sortedFiles[i].Status < sortedFiles[j].Status
		}
		return sortedFiles[i].Path < sortedFiles[j].Path
	})

	count := 0
	for _, f := range sortedFiles {
		if count >= 10 {
			s.WriteString(StyleDim.Render(fmt.Sprintf("... and %d more", len(sortedFiles)-10)))
			break
		}

		icon := " "
		color := StyleNormal

		switch f.Status {
		case "M":
			icon = "●"
			color = lipgloss.NewStyle().Foreground(ColorWarning)
		case "?":
			icon = "?"
			color = lipgloss.NewStyle().Foreground(ColorError)
		default:
			// check if staged
			if f.Staged {
				icon = "✓"
				color = lipgloss.NewStyle().Foreground(ColorSuccess)
			}
		}

		s.WriteString(fmt.Sprintf("%s %s\n", color.Render(icon), f.Path))
		count++
	}

	return StylePanel.Render(s.String())
}
