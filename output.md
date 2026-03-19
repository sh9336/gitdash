# GitDash

A terminal-based Git dashboard for visualizing repository information at a glance.

## Project Overview

GitDash provides a comprehensive terminal interface to monitor your Git repository status, including branches, recent commits, stash entries, and project statistics.

## Terminal Output Snapshot

```
 GitDash  
 
 • C/GoCodes/gitdash •  testbranch (HEAD moved)

 ╭────────────────────────────────────────────────────────────────╮
 │ ★ Branches (2)                                                 │
 │  ★ HEAD → testbranch (1 hour ago)                              │
 │  ▶ main (5 hours ago)                                          │
 ╰────────────────────────────────────────────────────────────────╯

 ╭────────────────────────────────────────────────────────────────╮
 │ Recent Commits                                                 │
 │ ────────────────────────────────────────────────────────────── │
 │ 63a6d6e feat: Branch Switch                                    │
 │         sh9336, 1 hour ago                                     │
 │ bbe4a53 style: Switch to Vertical Stack Layout for cleaner     │
 │         sh9336, 5 hours ago                                    │
 │ 25a1b41 style: Improve Help modal instructions (Press Esc)     │
 │         sh9336, 5 hours ago                                    │
 ╰────────────────────────────────────────────────────────────────╯

 ╭────────────────────────────────────────────────────────────────╮
 │ Stash                                                          │
 │ ────────────────────────────────────────────────────────────── │
 │   No stash entries                                             │
 ╰────────────────────────────────────────────────────────────────╯

 ╭────────────────────────────────────────────────────────────────╮
 │ Project Stats                                                  │
 │ ────────────────────────────────────────────────────────────── │
 │ Total Files: 16                                                │
 │ Go: 100.0%                                                     │
 ╰────────────────────────────────────────────────────────────────╯

 ╭────────────────────────────────────────────────────────────────╮
 │ Working Directory (On branch: testbranch)                      │
 │ ────────────────────────────────────────────────────────────── │
 │  ● internal/ui/dashboard.go                                    │
 │  ● internal/git/branches.go                                    │
 │  ● internal/git/repository.go                                  │
 ╰────────────────────────────────────────────────────────────────╯

 Press 'q' to quit, 'r' to refresh, '?' for help, 'Tab' to focus • Switched to branch: testbranch
```

## Technology Stack

- **Backend:** Go (100% of project)

## Controls

| Key | Action |
|-----|--------|
| `q` | Quit |
| `r` | Refresh |
| `?` | Help |
| `Tab` | Focus |

