package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitdash/gitdash/internal/config"
	"github.com/gitdash/gitdash/internal/git"
	"github.com/gitdash/gitdash/internal/stats"
)

type FocusArea int

const (
	FocusNone FocusArea = iota
	FocusBranches
)

type checkoutSuccessMsg struct {
	BranchName string
}

type checkoutTickMsg struct{}

type errMsg error

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
	Viewport      viewport.Model
	Quitting      bool
	Width         int
	Height        int
	Loading       bool
	ShowHelp      bool
	Focus         FocusArea
	StatusMessage string
	Spinner       int    // For checkout animation
	CheckingOut   string // Name of branch being checked out
}

func NewModel(info *git.RepoInfo, cfg *config.Config) Model {
	// Initial load
	m := Model{
		Config:      cfg,
		RepoInfo:    info,
		Loading:     true,
		ShowHelp:    false,
		Viewport:    viewport.New(0, 0),
		Focus:       FocusNone, // Start unfocused for normal dashboard scrolling
		Spinner:     0,
		CheckingOut: "",
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
	m.BranchesModel.Active = true // Since we default focus
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
			return errMsg(err)
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

func checkoutCmd(path string, branchName string, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		r, err := git.GetRepoInfo(path)
		if err != nil {
			return errMsg(err)
		}

		err = git.CheckoutBranch(r.Repo, branchName)
		if err != nil {
			return errMsg(err)
		}

		// Return success message - refresh will be handled separately
		return checkoutSuccessMsg{BranchName: branchName}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case checkoutTickMsg:
		if m.Loading {
			m.Spinner = (m.Spinner + 1) % 4
			m.Viewport.SetContent(m.RenderMainContent()) // Update viewport during animation
			return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return checkoutTickMsg{} })
		}

	case checkoutSuccessMsg:
		// Handle checkout success message
		m.StatusMessage = fmt.Sprintf("Checked out %s", msg.BranchName)
		m.Loading = false
		m.CheckingOut = "" // Clear checkout state
		// Trigger refresh with fresh repo info to get updated current branch
		freshRepoInfo, err := git.GetRepoInfo(m.RepoInfo.Path)
		if err != nil {
			m.StatusMessage = fmt.Sprintf("Error getting repo info: %v", err)
			return m, nil
		}
		// Update repo info immediately and trigger refresh command
		m.RepoInfo = freshRepoInfo
		return m, refreshData(freshRepoInfo, m.Config)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		headerHeight := 3
		footerHeight := 2
		verticalMarginHeight := headerHeight + footerHeight
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - verticalMarginHeight
		m.Viewport.SetContent(m.RenderMainContent())

	case refreshMsg:
		m.RepoInfo = msg.RepoInfo
		// Preserve selection if possible?
		// For now, reset is safer, or try to match name.
		// Simplest: just use new models.
		oldSelected := m.BranchesModel.Selected
		m.BranchesModel = msg.BranchesModel
		// Restore active state
		if m.Focus == FocusBranches {
			m.BranchesModel.Active = true
		}
		// Try to restore selection index if valid
		if oldSelected < len(m.BranchesModel.Branches) {
			// Actually NewBranchesModel tries to select current branch.
			// Users expects to stay on what they were hovering? Or jump to current?
			// If I just checked out, dragging selection to current is good.
			// So default behavior of NewBranchesModel (select current) is good for checkout.
		}

		m.CommitsModel = msg.CommitsModel
		m.WorkDirModel = msg.WorkDirModel
		m.StatsModel = msg.StatsModel
		m.Loading = false
		m.StatusMessage = fmt.Sprintf("Checked out %s", m.RepoInfo.CurrentBranch)
		m.Viewport.SetContent(m.RenderMainContent()) // Force viewport update

	case errMsg:
		m.StatusMessage = fmt.Sprintf("Error: %v", msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "r":
			m.Loading = true
			m.StatusMessage = "Refreshing..."
			return m, refreshData(m.RepoInfo, m.Config)
		case "?", "/":
			m.ShowHelp = !m.ShowHelp
			return m, nil
		case "esc":
			if m.ShowHelp {
				m.ShowHelp = false
				return m, nil
			}
			m.Quitting = true
			return m, tea.Quit
		case "tab":
			// Toggle focus
			if m.Focus == FocusNone {
				m.Focus = FocusBranches
				m.BranchesModel.Active = true
			} else {
				m.Focus = FocusNone
				m.BranchesModel.Active = false
			}
			m.Viewport.SetContent(m.RenderMainContent())
			return m, nil

		case "up", "k":
			if m.Focus == FocusBranches {
				m.BranchesModel.Previous()
				m.Viewport.SetContent(m.RenderMainContent())
				return m, nil
			}
		case "down", "j":
			if m.Focus == FocusBranches {
				m.BranchesModel.Next()
				m.Viewport.SetContent(m.RenderMainContent())
				return m, nil
			}
		case "enter":
			if m.Focus == FocusBranches {
				b := m.BranchesModel.Branches[m.BranchesModel.Selected]
				m.Loading = true
				m.StatusMessage = fmt.Sprintf("Checking out %s", b.Name)
				m.CheckingOut = b.Name
				m.Spinner = 0
				m.Viewport.SetContent(m.RenderMainContent()) // Update viewport immediately
				// Start spinner animation and checkout
				return m, tea.Batch(
					tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return checkoutTickMsg{} }),
					checkoutCmd(m.RepoInfo.Path, b.Name, m.Config),
				)
			}
		}
	}

	// Handle viewport scrolling if not focused on interactive element OR if we want to allow scrolling while branch is focused?
	// Usually arrow keys should do one thing.
	if m.Focus == FocusNone {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) RenderMainContent() string {
	panelWidth := m.Width - 4
	if panelWidth < 40 {
		panelWidth = 40
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"\n",
		m.BranchesModel.View(panelWidth, m.Loading, m.CheckingOut, m.Spinner),
		m.CommitsModel.View(panelWidth),
		m.StatsModel.View(panelWidth),
		m.WorkDirModel.View(panelWidth),
	)
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
	header := lipgloss.JoinVertical(lipgloss.Left,
		StyleTitle.Render("GitDash"),
		fmt.Sprintf(" • %s • %s", m.RepoInfo.Path, m.RepoInfo.CurrentBranch),
		"\n",
	)
	s.WriteString(header)

	s.WriteString(m.Viewport.View())

	// Footer
	spinnerChars := []string{"⠋", "⠙", "⠹", "⠸"}
	var spinner string
	if m.Loading {
		spinner = spinnerChars[m.Spinner] + " "
	}

	helpText := "Press 'q' to quit, 'r' to refresh, '?' for help, 'Tab' to focus"
	if m.Focus == FocusBranches {
		helpText += " • '↑/↓' select, 'Enter' checkout"
	} else {
		helpText += " • '↑/↓' to scroll"
	}

	if m.Loading {
		helpText += " • " + m.StatusMessage
	} else if m.StatusMessage != "" {
		helpText += " • " + m.StatusMessage
	}

	footer := StyleDim.Render("\n" + spinner + helpText)
	s.WriteString(footer)

	return s.String()
}
