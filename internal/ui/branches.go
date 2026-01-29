package ui

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gitdash/gitdash/internal/git"
)

type BranchesModel struct {
	Branches []git.Branch
	Selected int // For future navigation
}

func NewBranchesModel(branches []git.Branch) BranchesModel {
	return BranchesModel{
		Branches: branches,
	}
}

func (m BranchesModel) View() string {
	var s strings.Builder

	s.WriteString(StyleHeader.Render(fmt.Sprintf("Branches (%d)", len(m.Branches))))
	s.WriteString("\n")

	for i, b := range m.Branches {
		// Limit to showing top 5-10 branches for MVP
		if i >= 10 {
			s.WriteString(StyleDim.Render(fmt.Sprintf("... and %d more", len(m.Branches)-10)))
			break
		}

		cursor := "  " // Space for cursor
		if b.IsCurrent {
			cursor = "● "
		}

		nameStyle := StyleNormal
		if b.IsCurrent {
			nameStyle = StyleSelected
		}

		timeStr := humanize.Time(b.LastCommit)

		line := fmt.Sprintf("%s%s", cursor, nameStyle.Render(b.Name))
		// Right align time if possible, for now just append
		line += StyleDim.Render(fmt.Sprintf(" (%s)", timeStr))

		s.WriteString(line + "\n")
	}

	return StylePanel.Render(s.String())
}
