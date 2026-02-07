package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type FileStatus struct {
	Path   string
	Status string // M, A, D, ?, etc.
	Staged bool
}

type WorkingDirStatus struct {
	Files      []FileStatus
	BranchName string
	Modified   int
	Staged     int
	Untracked  int
	Conflicted int
}

// GetWorkingDirStatus returns the status of files in the working directory
func GetWorkingDirStatus(r *git.Repository) (*WorkingDirStatus, error) {
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := w.Status()
	if err != nil {
		return nil, err
	}

	// Determine if we are in detached HEAD or empty repo
	head, err := r.Head()
	branchName := "unknown"
	if err == nil {
		if head.Name().IsBranch() {
			branchName = head.Name().Short()
		} else {
			branchName = "Detached HEAD (" + head.Hash().String()[:7] + ")"
		}
	} else if err == plumbing.ErrReferenceNotFound {
		branchName = "Empty Repository"
	}

	ws := &WorkingDirStatus{
		Files:      []FileStatus{},
		BranchName: branchName,
	}

	for path, s := range status {
		code := string(s.Worktree)
		stagedCode := string(s.Staging)

		isStaged := stagedCode != " " && stagedCode != "?"
		isModified := code == "M" || stagedCode == "M"
		isUntracked := code == "?"
		isDeleted := code == "D" || stagedCode == "D"
		isAdded := stagedCode == "A"
		isConflicted := code == "C" || stagedCode == "C" || code == "U" || stagedCode == "U"

		// Determine display status
		displayStatus := code
		if isStaged {
			displayStatus = stagedCode
			ws.Staged++
		}

		if isModified && !isStaged {
			ws.Modified++
		}

		if isUntracked {
			displayStatus = "?"
			ws.Untracked++
		}

		if isConflicted {
			ws.Conflicted++
		}

		// Consider Deleted and Added in Modified/Staged counts?
		// For simplicity, let's just make sure they are counted.
		if isDeleted && !isStaged {
			ws.Modified++
		}

		if isAdded {
			// ws.Staged++ // Already handled by isStaged logic
		}

		ws.Files = append(ws.Files, FileStatus{
			Path:   path,
			Status: displayStatus,
			Staged: isStaged,
		})
	}

	// Sort files by path for consistent display?
	// Or maybe group by status (Staged, Modified, Untracked)
	// For MVP, just return list. UI can sort if needed.

	return ws, nil
}

// Helper to check if file is ignored?
// go-git w.Status() respects .gitignore automatically.
