package stats

import (
	"path/filepath"
	"sort"
	"strings"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type LanguageStat struct {
	Name       string
	Count      int
	Percentage float64
	Color      string
}

type ProjectStats struct {
	TotalFiles int
	Languages  []LanguageStat
	// Add other stats as needed
}

// Common extensions map
var extToLang = map[string]string{
	".go":   "Go",
	".js":   "JavaScript",
	".ts":   "TypeScript",
	".py":   "Python",
	".html": "HTML",
	".css":  "CSS",
	".md":   "Markdown",
	".yml":  "YAML",
	".yaml": "YAML",
	".json": "JSON",
	".c":    "C",
	".cpp":  "C++",
	".h":    "C",
	".java": "Java",
	".rs":   "Rust",
	// Add more as needed
}

// CalculateStats analyzes the repository content (simplified version)
func CalculateStats(r *gogit.Repository) (*ProjectStats, error) {
	// In a real implementation, we would walk the HEAD tree properly.
	// For simplicity and speed in MVP, we might just scan the worktree?
	// But standard practice is to scan HEAD.

	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	stats := &ProjectStats{
		Languages: []LanguageStat{},
	}

	langDates := make(map[string]int)

	// Scan the tree
	err = tree.Files().ForEach(func(f *object.File) error {
		stats.TotalFiles++

		ext := strings.ToLower(filepath.Ext(f.Name))
		if lang, ok := extToLang[ext]; ok {
			langDates[lang]++
		} else {
			// group unknown or no-extension files?
			// langDates["Other"]++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Calculate percentages
	var totalLangFiles int
	for _, count := range langDates {
		totalLangFiles += count
	}

	for name, count := range langDates {
		pct := 0.0
		if totalLangFiles > 0 {
			pct = (float64(count) / float64(totalLangFiles)) * 100
		}
		stats.Languages = append(stats.Languages, LanguageStat{
			Name:       name,
			Count:      count,
			Percentage: pct,
		})
	}

	// Sort by count desc
	sort.Slice(stats.Languages, func(i, j int) bool {
		return stats.Languages[i].Count > stats.Languages[j].Count
	})

	return stats, nil
}
