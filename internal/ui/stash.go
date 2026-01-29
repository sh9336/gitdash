package ui

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gitdash/gitdash/internal/git"
)

type StashModel struct {
	Entries []git.StashEntry
}

func NewStashModel(entries []git.StashEntry) StashModel {
	return StashModel{
		Entries: entries,
	}
}

func (m StashModel) View(width int) string {
	var s strings.Builder

	// Header
	s.WriteString(StyleHeader.Render("Stash"))
	s.WriteString("\n")

	if len(m.Entries) == 0 {
		s.WriteString(StyleDim.Render("   No stash entries"))
		return StylePanel.Copy().Width(width).Render(s.String())
	}

	for _, e := range m.Entries {
		id := fmt.Sprintf("stash@{%d}", e.ID)
		timeStr := humanize.Time(e.Timestamp)

		line := fmt.Sprintf(" %s %s %s",
			StyleSelected.Render(id),
			StyleNormal.Render(e.Message),
			StyleDim.Render("("+timeStr+")"),
		)
		s.WriteString(line + "\n")
	}

	return StylePanel.Copy().Width(width).Render(s.String())
}
