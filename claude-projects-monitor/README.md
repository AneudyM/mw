# claude-projects-monitor

A small command-line tool that inspects and watches the on-disk state Claude
Code keeps under `~/.claude/projects/`. Each subdirectory there represents a
working directory Claude has been used in, and each `*.jsonl` file inside is a
session transcript.

## Build

```
go build -o cpm ./cmd/cpm
```

## Usage

```
cpm list                 # list projects with session counts and last activity
cpm sessions <project>   # list sessions inside one project
cpm stats                # aggregate stats across all projects
cpm watch                # poll the projects dir and report new activity
```

`<project>` may be either the encoded directory name (as stored on disk) or
the decoded working-directory path. Use `--root <path>` on any command to
point at a different projects root.

## Layout

```
cmd/cpm/        CLI entry point
internal/projects/   project + session discovery and parsing
internal/monitor/    poll-based change watcher
internal/stats/      aggregation helpers
```

No third-party dependencies — only the Go standard library.
