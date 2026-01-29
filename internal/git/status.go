package git

import (
	"github.com/go-git/go-git/v5"
)

type FileStatus struct {
	Path   string
	Status string // M, A, D, ?, etc.
	Staged bool
}

type WorkingDirStatus struct {
	Files      []FileStatus
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

	ws := &WorkingDirStatus{
		Files: []FileStatus{},
	}

	for path, s := range status {
		code := string(s.Worktree)
		stagedCode := string(s.Staging)

		isStaged := stagedCode != " " && stagedCode != "?"
		isModified := code == "M"
		isUntracked := code == "?"

		// Determine display status
		displayStatus := code
		if isStaged {
			displayStatus = stagedCode
			ws.Staged++
		}

		if isModified {
			ws.Modified++
		}

		if isUntracked {
			displayStatus = "?"
			ws.Untracked++
		}

		if code == "C" || stagedCode == "C" {
			ws.Conflicted++
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
