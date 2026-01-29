package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitdash/gitdash/internal/git"
	"github.com/gitdash/gitdash/internal/stats"
)

type Model struct {
	RepoInfo      *git.RepoInfo
	BranchesModel BranchesModel
	CommitsModel  CommitsModel
	WorkDirModel  WorkDirModel
	StatsModel    StatsModel
	Quitting      bool
	Width         int
	Height        int
}

func NewModel(info *git.RepoInfo) Model {
	// Initialize sub-models
	branches, _ := git.GetBranches(info.Repo)
	commits, _ := git.GetRecentCommits(info.Repo, 10)
	status, _ := git.GetWorkingDirStatus(info.Repo)
	projectStats, _ := stats.CalculateStats(info.Repo)

	return Model{
		RepoInfo:      info,
		BranchesModel: NewBranchesModel(branches),
		CommitsModel:  NewCommitsModel(commits),
		WorkDirModel:  NewWorkDirModel(status),
		StatsModel:    NewStatsModel(projectStats),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var s strings.Builder

	// Header
	s.WriteString(StyleTitle.Render("GitDash"))
	s.WriteString(fmt.Sprintf(" • %s • %s", m.RepoInfo.Path, m.RepoInfo.CurrentBranch))
	s.WriteString("\n\n")

	// Main Layout
	// Row 1: Branches | Commits
	row1 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.BranchesModel.View(),
		m.CommitsModel.View(),
	)

	// Row 2: WorkDir | Stats
	row2 := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.WorkDirModel.View(),
		m.StatsModel.View(),
	)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		row1,
		row2,
	)

	s.WriteString(mainContent)

	// Help
	s.WriteString(StyleDim.Render("\nPress 'q' to quit"))

	return lipgloss.Place(m.Width, m.Height, lipgloss.Top, lipgloss.Left, s.String())
}
