package git

import (
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Branch struct {
	Name       string
	IsCurrent  bool
	LastCommit time.Time
	Hash       string
	Remote     string
	RemoteIdx  int // For handling multiple remotes if needed, simplified here
}

// GetBranches returns a list of local branches sorted by recency
func GetBranches(r *git.Repository) ([]Branch, error) {
	branches := []Branch{}

	// Get HEAD to check current branch
	headRef, err := r.Head()
	currentBranchName := ""
	if err == nil {
		currentBranchName = headRef.Name().Short()
	}

	bs, err := r.Branches()
	if err != nil {
		return nil, err
	}

	err = bs.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		isCurrent := name == currentBranchName

		// Get last commit for this branch
		commit, err := r.CommitObject(ref.Hash())
		var lastCommit time.Time
		if err == nil {
			lastCommit = commit.Author.When
		}

		// Basic remote tracking logic (simplified)
		// To get real tracking info requires config parsing, for now just placeholder

		branches = append(branches, Branch{
			Name:       name,
			IsCurrent:  isCurrent,
			LastCommit: lastCommit,
			Hash:       ref.Hash().String(),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort by last commit date (descending)
	sort.Slice(branches, func(i, j int) bool {
		return branches[i].LastCommit.After(branches[j].LastCommit)
	})

	return branches, nil
}

// CheckoutBranch checks out the given branch name
func CheckoutBranch(r *git.Repository, branchName string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
	})
}
