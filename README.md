# tasks

A CLI and TUI TODO list manager with a cyberpunk-themed interface.

*Inspired by John Calhoun's Gophercises*

## Features

- **Interactive TUI** – Run without arguments for a full terminal UI
- **CLI mode** – Use commands for scripting and automation
- **Buckets** – Organize tasks in separate lists (e.g. work, personal)
- **Priorities** – alta (high), media (medium), baixa (low)
- **Due dates** – Set and filter by due date
- **Archive** – Mark tasks complete without deleting; view archived tasks
- **Search/filter** – Filter by text or priority
- **Pomodoro timer** – Built-in 25-minute focus timer

## Installation

```bash
git clone https://github.com/lmagdanello/tasks.git
cd tasks
go install .
```

Or download a pre-built binary from the [GitHub Actions](https://github.com/lmagdanello/tasks/actions) artifacts.

## Quick Start

```bash
# Launch the interactive TUI (default)
tasks

# Or explicitly
tasks --tui

# Use CLI commands
tasks init --bucket mylist
tasks add --bucket mylist "First task"
tasks list --bucket mylist
```

## TUI Mode

Run `tasks` with no arguments to open the TUI. Navigate with **Tab** / **Shift+Tab**.

| Key     | Action        |
|---------|---------------|
| `Tab`   | Navigate      |
| `Enter` | Task menu     |
| `a`     | Add task      |
| `e`     | Edit task     |
| `d`     | Delete task   |
| `n`     | New bucket    |
| `/`     | Filter        |
| `v`     | Toggle archived view |
| `p`     | Pomodoro start/stop   |
| `q`     | Quit          |

## CLI Commands

| Command | Description |
|---------|-------------|
| `init --bucket <name>` | Create a new bucket |
| `add -b <bucket> "task" [-p priority] [-d YYYY-MM-DD]` | Add a task |
| `list -b <bucket> [--done] [-f filter] [-p priority]` | List tasks |
| `edit -b <bucket> <id> [text] [-p priority] [-d due]` | Edit a task |
| `done -b <bucket> <id> [--delete]` | Archive (or delete) a task |
| `search -b <bucket> [query]` | Search tasks by text |

### Examples

```bash
# Create a bucket
tasks init --bucket github

# Add tasks with options
tasks add -b github "Fix bug" -p alta -d 2025-02-10
tasks add -b github "Write docs"

# List and filter
tasks list -b github
tasks list -b github --done
tasks list -b github -f "bug" -p alta

# Edit and archive
tasks edit -b github 1 "Fix critical bug"
tasks done -b github 1
tasks done -b github 2 --delete  # permanent delete
```

## Storage

Tasks are stored in `~/.tasks/`. Each bucket is a `.db` file (e.g. `~/.tasks/github.db`). Data is persisted using [bbolt](https://github.com/etcd-io/bbolt).

## Contributing

Contributions are welcome. Feel free to open issues or pull requests.
