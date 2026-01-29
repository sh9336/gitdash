package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRepo(t *testing.T) {
	// Create a temp dir
	tmpDir, err := os.MkdirTemp("", "gitdash-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a .git dir inside
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a subdir
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Test finding from root
	// On Windows, paths can be tricky with case/separator, so we normalize for comparison if needed,
	// but FindRepo uses filepath.Abs so it should be consistent.
	// However, MkdirTemp returns absolute path on Windows usually?

	found, err := FindRepo(tmpDir)
	if err != nil {
		t.Errorf("FindRepo(tmpDir) returned error: %v", err)
	}

	// Normalize paths for comparison (Resolve symlinks if any, though MkdirTemp usually real path)
	// Just string compare should work if implementation returns filepath.Abs

	if found != tmpDir {
		t.Errorf("FindRepo(tmpDir) = %v; want %v", found, tmpDir)
	}

	// Test finding from subdir
	found, err = FindRepo(subDir)
	if err != nil {
		t.Errorf("FindRepo(subDir) returned error: %v", err)
	}
	if found != tmpDir { // Should find parent
		t.Errorf("FindRepo(subDir) = %v; want %v", found, tmpDir)
	}

	// Test non-git dir
	nonGitDir, err := os.MkdirTemp("", "gitdash-nongit")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(nonGitDir)

	_, err = FindRepo(nonGitDir)
	if err == nil {
		t.Error("FindRepo(nonGitDir) should fail")
	}
}
