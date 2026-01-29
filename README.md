# GitDash

GitDash is a terminal UI dashboard that provides developers with an instant, visual overview of their Git repositories. Built with Go and Bubble Tea.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)

## Features

- **📊 Dashboard Overview**: View branch status, recent commits, and working directory state in one view.
- **🌳 Branch Management**: List local branches with recency and tracking status.
- **📝 Commit History**: Recent commits with hash, author, and relative time.
- **📁 File Status**: Visual indicators for modified, staged, and untracked files.
- **📈 Project Stats**: Language distribution and file counts.

## Installation

### From Source

```bash
git clone https://github.com/yourusername/gitdash.git
cd gitdash
go install cmd/gitdash/main.go
# or
make build
```

## Usage

Run GitDash in your git repository:

```bash
gitdash
# or specify a path
gitdash --path /path/to/repo
```

### Controls

- `q` or `Ctrl+C`: Quit application
- `r`: Refresh data (Coming soon)

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Go-Git](https://github.com/go-git/go-git)
- [Cobra](https://github.com/spf13/cobra)

## License

MIT
