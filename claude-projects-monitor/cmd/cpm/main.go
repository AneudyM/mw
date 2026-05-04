// Command cpm — Claude Projects Monitor — inspects and watches the on-disk
// state Claude Code keeps under ~/.claude/projects.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/aneudym/mw/claude-projects-monitor/internal/monitor"
	"github.com/aneudym/mw/claude-projects-monitor/internal/projects"
	"github.com/aneudym/mw/claude-projects-monitor/internal/stats"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "cpm:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		usage(stderr)
		return fmt.Errorf("missing command")
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "list":
		return cmdList(rest, stdout)
	case "sessions":
		return cmdSessions(rest, stdout)
	case "stats":
		return cmdStats(rest, stdout)
	case "watch":
		return cmdWatch(rest, stdout)
	case "help", "-h", "--help":
		usage(stdout)
		return nil
	default:
		usage(stderr)
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `cpm — Claude Projects Monitor

Usage:
  cpm list                 list projects with session counts and last activity
  cpm sessions <project>   list sessions inside one project
  cpm stats                aggregate stats across all projects
  cpm watch                poll the projects dir and report new activity

Common flags:
  --root <path>            override the projects root (default: ~/.claude/projects)
`)
}

func resolveRoot(root string) (string, error) {
	if root != "" {
		return root, nil
	}
	return projects.DefaultRoot()
}

func addRootFlag(fs *flag.FlagSet, root *string) {
	fs.StringVar(root, "root", "", "projects root (default: ~/.claude/projects)")
}

func cmdList(args []string, w io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	var root string
	addRootFlag(fs, &root)
	if err := fs.Parse(args); err != nil {
		return err
	}
	r, err := resolveRoot(root)
	if err != nil {
		return err
	}
	ps, err := projects.List(r)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROJECT\tSESSIONS\tSIZE\tLAST ACTIVITY")
	for _, p := range ps {
		last := "-"
		if t := p.LastActivity(); !t.IsZero() {
			last = t.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n", p.EncodedName, len(p.Sessions), humanBytes(p.TotalSize()), last)
	}
	return tw.Flush()
}

func cmdSessions(args []string, w io.Writer) error {
	fs := flag.NewFlagSet("sessions", flag.ContinueOnError)
	var root string
	addRootFlag(fs, &root)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("sessions: project identifier required")
	}
	r, err := resolveRoot(root)
	if err != nil {
		return err
	}
	p, err := projects.Find(r, fs.Arg(0))
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Project: %s\nPath:    %s\nDecoded: %s\n\n", p.EncodedName, p.Path, p.DecodedDir())
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SESSION\tSIZE\tENTRIES\tMODIFIED")
	for _, s := range p.Sessions {
		entries, err := projects.CountEntries(s.Path)
		if err != nil {
			return err
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n", s.ID, humanBytes(s.Size), entries, s.Modified.Format(time.RFC3339))
	}
	return tw.Flush()
}

func cmdStats(args []string, w io.Writer) error {
	fs := flag.NewFlagSet("stats", flag.ContinueOnError)
	var root string
	var top int
	addRootFlag(fs, &root)
	fs.IntVar(&top, "top", 5, "show top-N projects by session count")
	if err := fs.Parse(args); err != nil {
		return err
	}
	r, err := resolveRoot(root)
	if err != nil {
		return err
	}
	ps, err := projects.List(r)
	if err != nil {
		return err
	}
	s := stats.Aggregate(ps, top)
	last := "-"
	if !s.LastActivity.IsZero() {
		last = s.LastActivity.Format(time.RFC3339)
	}
	fmt.Fprintf(w, "Root:          %s\n", r)
	fmt.Fprintf(w, "Projects:      %d\n", s.Projects)
	fmt.Fprintf(w, "Sessions:      %d\n", s.Sessions)
	fmt.Fprintf(w, "Total size:    %s\n", humanBytes(s.TotalBytes))
	fmt.Fprintf(w, "Last activity: %s\n", last)
	if len(s.Top) == 0 {
		return nil
	}
	fmt.Fprintf(w, "\nTop %d by session count:\n", len(s.Top))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "  PROJECT\tSESSIONS\tSIZE\tLAST")
	for _, t := range s.Top {
		l := "-"
		if !t.Last.IsZero() {
			l = t.Last.Format(time.RFC3339)
		}
		fmt.Fprintf(tw, "  %s\t%d\t%s\t%s\n", t.Name, t.Sessions, humanBytes(t.Bytes), l)
	}
	return tw.Flush()
}

func cmdWatch(args []string, w io.Writer) error {
	fs := flag.NewFlagSet("watch", flag.ContinueOnError)
	var root string
	var interval time.Duration
	addRootFlag(fs, &root)
	fs.DurationVar(&interval, "interval", 2*time.Second, "poll interval")
	if err := fs.Parse(args); err != nil {
		return err
	}
	r, err := resolveRoot(root)
	if err != nil {
		return err
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	fmt.Fprintf(w, "watching %s (interval %s) — ^C to stop\n", r, interval)
	events, errs := monitor.Watch(ctx, r, interval)
	for {
		select {
		case <-ctx.Done():
			return nil
		case e, ok := <-events:
			if !ok {
				return nil
			}
			fmt.Fprintf(w, "%s  %-7s  %s/%s  %s\n",
				e.When.Format(time.RFC3339), e.Kind, shortProject(e.Project), e.Session, humanBytes(e.Size))
		case err, ok := <-errs:
			if !ok {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}
}

func shortProject(name string) string {
	if len(name) <= 32 {
		return name
	}
	return "…" + name[len(name)-31:]
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return strings.TrimSuffix(fmt.Sprintf("%.1f", float64(n)/float64(div)), ".0") +
		string("KMGTPE"[exp]) + "B"
}
