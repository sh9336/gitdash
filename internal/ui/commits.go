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

func (m CommitsModel) View() string {
	var s strings.Builder

	s.WriteString(StyleHeader.Render("Recent Commits"))
	s.WriteString("\n")

	if len(m.Commits) == 0 {
		s.WriteString(StyleDim.Render("  No commits found"))
		return StylePanel.Render(s.String())
	}

	for _, c := range m.Commits {
		// Truncate hash
		hash := c.Hash[:7]

		// Truncate message if too long (simple approach)
		msg := strings.Split(c.Message, "\n")[0] // First line only
		if len(msg) > 50 {
			msg = msg[:47] + "..."
		}

		timeStr := humanize.Time(c.When)

		// Layout: Hash Msg Author Time
		// a7f3d92 feat: add stats panel
		//         John Doe, 2 hours ago

		line1 := fmt.Sprintf("%s %s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(hash),
			StyleNormal.Render(msg),
		)

		line2 := fmt.Sprintf("        %s, %s",
			StyleDim.Render(c.Author),
			StyleDim.Render(timeStr),
		)

		s.WriteString(line1 + "\n" + line2 + "\n")
	}

	return StylePanel.Render(s.String())
}
