package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Commit struct {
	Hash        string
	Message     string
	Author      string
	AuthorEmail string
	When        time.Time
}

// GetRecentCommits returns the last n commits from the given branch (or HEAD if empty)
func GetRecentCommits(r *git.Repository, branchName string, n int) ([]Commit, error) {
	var hash plumbing.Hash

	if branchName == "" {
		ref, err := r.Head()
		if err != nil {
			return nil, nil // Empty repo or other error, return empty list
		}
		hash = ref.Hash()
	} else {
		ref, err := r.Reference(plumbing.ReferenceName("refs/heads/"+branchName), true)
		if err != nil {
			return nil, err
		}
		hash = ref.Hash()
	}

	cIter, err := r.Log(&git.LogOptions{From: hash})
	if err != nil {
		return nil, err
	}

	var commits []Commit
	count := 0
	err = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, Commit{
			Hash:        c.Hash.String(),
			Message:     c.Message,
			Author:      c.Author.Name,
			AuthorEmail: c.Author.Email,
			When:        c.Author.When,
		})
		count++

		if count >= n {
			return storer.ErrStop
		}
		return nil
	})

	return commits, nil
}
