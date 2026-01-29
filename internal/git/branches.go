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
	currentHash := ""
	if err == nil {
		currentHash = headRef.Hash().String()
	}

	bs, err := r.Branches()
	if err != nil {
		return nil, err
	}

	err = bs.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		isCurrent := ref.Hash().String() == currentHash

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
