package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/gitdash/gitdash/internal/git"
)

type CommitsModel struct {
	Commits []git.Commit
}

func NewCommitsModel(commits []git.Commit) CommitsModel {
	return CommitsModel{
		Commits: commits,
	}
}

func (m CommitsModel) View(width int) string {
	var s strings.Builder

	// Header
	s.WriteString(StyleHeader.Render("Recent Commits"))
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", width)))
	s.WriteString("\n")

	if len(m.Commits) == 0 {
		s.WriteString(StyleDim.Render("   No commits found"))
		return StylePanel.Copy().Width(width).Render(s.String())
	}

	for _, c := range m.Commits {
		// Truncate hash
		hash := c.Hash[:7]

		// Message truncation
		msg := strings.Split(c.Message, "\n")[0]
		maxLen := width - 14 // 14 gives space for padding check
		if maxLen < 10 {
			maxLen = 10
		}
		if len(msg) > maxLen {
			msg = msg[:maxLen-3] + "..."
		}

		timeStr := humanize.Time(c.When)

		// Layout: Hash Msg Author Time
		// Add padding " "
		line1 := fmt.Sprintf(" %s %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(hash),
			StyleNormal.Render(msg),
		)

		line2 := fmt.Sprintf("         %s, %s",
			StyleDim.Render(c.Author),
			StyleDim.Render(timeStr),
		)

		s.WriteString(line1 + "\n" + line2 + "\n")
	}

	return StylePanel.Copy().Width(width).Render(s.String())
}
