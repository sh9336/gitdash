package git

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type RepoInfo struct {
	Path          string
	CurrentBranch string
	IsClean       bool
	Remotes       []string
	Repo          *git.Repository
}

// FindRepo walks up the directory tree to find a .git directory
func FindRepo(startPath string) (string, error) {
	path, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
			return path, nil
		}

		parent := filepath.Dir(path)
		if parent == path {
			return "", errors.New("git repository not found")
		}
		path = parent
	}
}

// GetRepoInfo opens the repository and extracts basic information
func GetRepoInfo(path string) (*RepoInfo, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Get current branch
	headRef, err := r.Head()
	currentBranch := ""
	if err == nil {
		if headRef.Name().IsBranch() {
			currentBranch = headRef.Name().Short()
		} else {
			currentBranch = headRef.Hash().String()[:7] // Detached head
		}
	} else if err == plumbing.ErrReferenceNotFound {
		currentBranch = "Empty Repository"
	}

	// Check status
	isClean := true
	w, err := r.Worktree()
	if err == nil {
		status, err := w.Status()
		if err == nil && !status.IsClean() {
			isClean = false
		}
	}

	// Get remotes
	remotes, err := r.Remotes()
	var remoteURLs []string
	if err == nil {
		for _, rem := range remotes {
			if len(rem.Config().URLs) > 0 {
				remoteURLs = append(remoteURLs, rem.Config().URLs...)
			}
		}
	}

	return &RepoInfo{
		Path:          path,
		CurrentBranch: currentBranch,
		IsClean:       isClean,
		Remotes:       remoteURLs,
		Repo:          r,
	}, nil
}
