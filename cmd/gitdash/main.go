package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sh9336/gitdash/internal/config"
	"github.com/sh9336/gitdash/internal/git"
	"github.com/sh9336/gitdash/internal/ui"
	"github.com/spf13/cobra"
)

var (
	pathFlag   string
	configFlag string
)

const Version = "1.0.1"

func main() {
	var rootCmd = &cobra.Command{
		Use:     "gitdash",
		Short:   "GitDash is a terminal UI dashboard for git repositories",
		Version: Version,
		Run:     run,
	}

	rootCmd.PersistentFlags().StringVarP(&pathFlag, "path", "p", ".", "Path to git repository")
	rootCmd.PersistentFlags().StringVarP(&configFlag, "config", "c", "", "Path to configuration file")

	// Set version template
	rootCmd.SetVersionTemplate("GitDash version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Load config
	cfg, err := config.LoadConfig(configFlag)
	if err != nil {
		// Just warn or ignore for now, use defaults
		fmt.Printf("Warning: could not load config: %v\n", err)
		// continue with nil config (defaults will be used in NewModel logic if nil check exists, but easier to just use default struct)
		// But LoadConfig returns error if viper fails significantly.
		// If just file not found, we ignored it in LoadConfig but maybe it returns nil?
		// LoadConfig returns &Config with defaults if file missing.
	}

	repoPath, err := git.FindRepo(pathFlag)
	if err != nil {
		fmt.Printf("Error finding repository: %v\n", err)
		os.Exit(1)
	}

	info, err := git.GetRepoInfo(repoPath)
	if err != nil {
		fmt.Printf("Error getting repository info: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.NewModel(info, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running dashboard: %v\n", err)
		os.Exit(1)
	}
}
