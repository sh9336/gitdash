# 🚀 GitDash

**GitDash** is a fast, beautiful, and non-destructive terminal dashboard for your Git repositories. It provides a real-time overview of your branches, commits, stashes, and project statistics without putting your workspace at risk.

![GitDash Preview](output.md)

## 🌟 Key Features

- **🔍 Safe Inspection Mode**: Fly through your branches with arrow keys. GitDash automatically fetches history and stats for the selected branch *without* checking it out physically. Your uncommitted work is 100% safe.
- **📊 Real-time Analytics**: See project language composition, file counts, and commit velocity at a glance.
- **🛡️ Workspace Awareness**: Clear visibility of your working directory status (Modified, Staged, Untracked, Conflicted).
- **⌨️ Keyboard Centric**: Designed for speed with intuitive Vim-style navigation.
- **🎨 Premium Aesthetics**: Built with `BubbleTea` and `LipGloss` for a stunning terminal experience.

## 💾 Installation

Ensure you have [Go](https://golang.org/dl/) installed (version 1.18+), then run:

```bash
go install github.com/gitdash/gitdash/cmd/gitdash@latest
```

Ensure your `GOBIN` directory (typically `~/go/bin`) is in your system `PATH`.

## 🚀 Quick Start

Launch GitDash from any git repository:

```bash
cd /your/git/repo
gitdash
```

Alternatively, specify a path:

```bash
gitdash --path /path/to/repo
```

## ⌨️ Controls

| Key | Action |
|-----|--------|
| `Tab` | Toggle Focus between main scroll and Branches list |
| `↑ / ↓` | Scroll dashboard OR **Inspect** selected branch |
| `f` | **Force Checkout** (Discards local changes to switch) |
| `r` | Hard Refresh all data |
| `?` | Toggle Help modal |
| `q / Esc` | Quit GitDash |

## ⚙️ Configuration

GitDash looks for a `.gitdash.yaml` in your project root or home directory.

```yaml
commits:
  show_count: 15
  show_author: true
  show_relative_time: true

display:
  colors: true
  unicode: true

dashboard:
  refresh_interval: "30s"
```

## 🛠️ Performance

GitDash is built for speed. It uses the `go-git` library for direct object access, meaning it doesn't need to shell out to the `git` binary for every update. It utilizes `storer.ErrStop` optimization to ensure large histories are truncated and loaded efficiently.

## 🎯 Use Case Matrix

| Use Case | How GitDash Helps |
|----------|-------------------|
| **Code Review** | Inspect a PR branch's commits instantly without checking it out locally. |
| **Repo Auditing** | View language distribution and file counts across different versions of the code. |
| **Workspace Safety** | Browse the repo while keeping your complicated, uncommitted setup untouched. |
| **Snappy Navigation** | Fast, Vim-style scrolling through large commit histories. |

## 🤝 Contributing

Contributions are welcome! Whether it's a bug report, feature request, or a Pull Request:

1. Fork the repo.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.
