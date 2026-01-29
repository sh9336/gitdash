package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitdash/gitdash/internal/config"
	"github.com/gitdash/gitdash/internal/git"
	"github.com/gitdash/gitdash/internal/stats"
)

type refreshMsg struct {
	RepoInfo      *git.RepoInfo
	BranchesModel BranchesModel
	CommitsModel  CommitsModel
	WorkDirModel  WorkDirModel
	StatsModel    StatsModel
}

type Model struct {
	Config        *config.Config
	RepoInfo      *git.RepoInfo
	BranchesModel BranchesModel
	CommitsModel  CommitsModel
	WorkDirModel  WorkDirModel
	StatsModel    StatsModel
	Quitting      bool
	Width         int
	Height        int
	Loading       bool
	ShowHelp      bool
}

func NewModel(info *git.RepoInfo, cfg *config.Config) Model {
	// Initial load
	m := Model{
		Config:   cfg,
		RepoInfo: info,
		Loading:  true,
		ShowHelp: false,
	}

	// Use config for commit count
	commitCount := 10
	if cfg != nil && cfg.Commits.ShowCount > 0 {
		commitCount = cfg.Commits.ShowCount
	}

	branches, _ := git.GetBranches(info.Repo)
	commits, _ := git.GetRecentCommits(info.Repo, commitCount)
	status, _ := git.GetWorkingDirStatus(info.Repo)
	projectStats, _ := stats.CalculateStats(info.Repo)

	m.BranchesModel = NewBranchesModel(branches)
	m.CommitsModel = NewCommitsModel(commits)
	m.WorkDirModel = NewWorkDirModel(status)
	m.StatsModel = NewStatsModel(projectStats)
	m.Loading = false

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func refreshData(info *git.RepoInfo, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		newInfo, err := git.GetRepoInfo(info.Path)
		if err != nil {
			return nil
		}

		commitCount := 10
		if cfg != nil && cfg.Commits.ShowCount > 0 {
			commitCount = cfg.Commits.ShowCount
		}

		branches, _ := git.GetBranches(newInfo.Repo)
		commits, _ := git.GetRecentCommits(newInfo.Repo, commitCount)
		status, _ := git.GetWorkingDirStatus(newInfo.Repo)
		projectStats, _ := stats.CalculateStats(newInfo.Repo)

		return refreshMsg{
			RepoInfo:      newInfo,
			BranchesModel: NewBranchesModel(branches),
			CommitsModel:  NewCommitsModel(commits),
			WorkDirModel:  NewWorkDirModel(status),
			StatsModel:    NewStatsModel(projectStats),
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case refreshMsg:
		m.RepoInfo = msg.RepoInfo
		m.BranchesModel = msg.BranchesModel
		m.CommitsModel = msg.CommitsModel
		m.WorkDirModel = msg.WorkDirModel
		m.StatsModel = msg.StatsModel
		m.Loading = false

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "r":
			m.Loading = true
			return m, refreshData(m.RepoInfo, m.Config)
		case "?", "/": // Handle both ? and / for help
			m.ShowHelp = !m.ShowHelp
			return m, nil
		case "esc":
			if m.ShowHelp {
				m.ShowHelp = false
				return m, nil
			}
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

	if m.ShowHelp {
		return m.helpView()
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
	helpText := "Press 'q' to quit, 'r' to refresh, '?' for help"
	if m.Loading {
		helpText += " • Refreshing..."
	}
	s.WriteString(StyleDim.Render("\n" + helpText))

	return lipgloss.Place(m.Width, m.Height, lipgloss.Top, lipgloss.Left, s.String())
}
