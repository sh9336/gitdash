package ui

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/gitdash/gitdash/internal/git"
)

type BranchesModel struct {
	Branches []git.Branch
	Selected int
	Active   bool // Whether this panel is currently active/focused
}

func NewBranchesModel(branches []git.Branch) BranchesModel {
	m := BranchesModel{
		Branches: branches,
		Selected: 0,
		Active:   false,
	}
	// Try to select the current branch by default
	for i, b := range branches {
		if b.IsCurrent {
			m.Selected = i
			break
		}
	}
	return m
}

func (m *BranchesModel) Next() {
	if m.Selected < len(m.Branches)-1 {
		m.Selected++
	}
}

func (m *BranchesModel) Previous() {
	if m.Selected > 0 {
		m.Selected--
	}
}

func (m BranchesModel) View(width int, loading bool, checkingOut string, spinner int) string {
	var s strings.Builder

	// Header
	title := fmt.Sprintf("Branches (%d)", len(m.Branches))
	if m.Active {
		title = "★ " + title // Indicate focus
		s.WriteString(StyleSelected.Copy().Bold(true).Render(title))
	} else {
		s.WriteString(StyleHeader.Render(title))
	}
	s.WriteString("\n")

	// Limit number of branches shown to fit in panel roughly, or just show all (viewport handles scrolling main area)
	// But wait, the Dashboard Viewport handles the *whole* content.
	// If the list is huge, we might need internal paging, but for now let's rely on the main viewport or just show top N.
	// Actually, if we are navigating *inside* this panel, we need to ensure the selected item is visible if we implement scrolling *inside* the panel.
	// For now, let's just render all and let the user scroll the whole dashboard?
	// No, if I press Down, I expect the selection to move.

	// Let's render all for now.
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸"}
	for i, b := range m.Branches {
		// Add spinner if this branch is being checked out
		spinnerPrefix := ""
		if loading && checkingOut == b.Name {
			spinnerPrefix = spinnerChars[spinner%4] + " "
		}

		cursor := "  "
		if b.IsCurrent {
			cursor = " ●" // Bullet for current branch
		}

		// Visual indicator for selection cursor
		if m.Active && i == m.Selected {
			cursor = " >"
			if b.IsCurrent {
				cursor = " >●" // Both selection and current branch
			}
		}

		nameStyle := StyleNormal
		if b.IsCurrent {
			nameStyle = StyleSelected // Current branch in green/bold
		}

		// Highlight the entire line if selected and active
		lineContent := b.Name
		if b.IsCurrent && !m.Active {
			lineContent = "★ " + b.Name // Star indicator for current branch when unfocused
		}

		timeStr := humanize.Time(b.LastCommit)
		rightContent := fmt.Sprintf(" (%s)", timeStr)

		mainText := nameStyle.Render(lineContent)
		if m.Active && i == m.Selected {
			mainText = StyleSelected.Copy().Underline(true).Render(lineContent)
		}

		line := fmt.Sprintf(" %s%s %s", spinnerPrefix, cursor, mainText)
		line += StyleDim.Render(rightContent)

		s.WriteString(line + "\n")
	}

	style := StylePanel.Copy().Width(width)
	if m.Active {
		style = style.BorderForeground(ColorPrimary)
	}

	return style.Render(s.String())
}
