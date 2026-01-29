package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Commit struct {
	Hash        string
	Message     string
	Author      string
	AuthorEmail string
	When        time.Time
}

// GetRecentCommits returns the last n commits from the current HEAD
func GetRecentCommits(r *git.Repository, n int) ([]Commit, error) {
	// check if head exists (might be empty repo)
	ref, err := r.Head()
	if err != nil {
		return nil, nil // Empty repo or other error, return empty list
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}

	var commits []Commit
	count := 0
	err = cIter.ForEach(func(c *object.Commit) error {
		if count >= n {
			return nil // Stop iterating
		}

		commits = append(commits, Commit{
			Hash:        c.Hash.String(),
			Message:     c.Message,
			Author:      c.Author.Name,
			AuthorEmail: c.Author.Email,
			When:        c.Author.When,
		})
		count++
		return nil
	})

	return commits, nil
}
