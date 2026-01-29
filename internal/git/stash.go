package git

import (
	"time"

	"github.com/go-git/go-git/v5"
)

type StashEntry struct {
	ID        int
	Message   string
	Timestamp time.Time
	Hash      string
}

// GetStashList returns the list of stash entries.
// Note: go-git support for stash is limited, we might need to parse refs/stash.
func GetStashList(r *git.Repository) ([]StashEntry, error) {
	// For now, returning empty to implement UI structure.
	// Real implementation would iterate over refs/stash and its reflog.
	return []StashEntry{}, nil
}
