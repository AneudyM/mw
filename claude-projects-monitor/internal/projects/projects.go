// Package projects discovers Claude Code projects on disk and the session
// transcripts stored inside each one.
package projects

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DefaultRoot returns the conventional projects directory, ~/.claude/projects.
func DefaultRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", "projects"), nil
}

// Project is one directory under the projects root.
type Project struct {
	// EncodedName is the on-disk directory name. Claude Code derives it from
	// the working directory by replacing path separators with '-', so
	// "/Users/foo/code" becomes "-Users-foo-code".
	EncodedName string
	// Path is the absolute path to the project directory.
	Path string
	// Sessions are the .jsonl transcript files inside the project.
	Sessions []Session
}

// Session is one transcript file inside a project.
type Session struct {
	ID       string // file name without the .jsonl suffix
	Path     string
	Size     int64
	Modified time.Time
}

// DecodedDir returns a best-effort decoding of EncodedName back to a path.
// The encoding is lossy — '-' inside real path components becomes
// indistinguishable from a separator — so the result is informational only.
func (p Project) DecodedDir() string {
	if p.EncodedName == "" {
		return ""
	}
	return strings.ReplaceAll(p.EncodedName, "-", "/")
}

// LastActivity is the most recent session modification time, or the zero
// value if the project has no sessions.
func (p Project) LastActivity() time.Time {
	var t time.Time
	for _, s := range p.Sessions {
		if s.Modified.After(t) {
			t = s.Modified
		}
	}
	return t
}

// TotalSize sums the sizes of every session in the project.
func (p Project) TotalSize() int64 {
	var n int64
	for _, s := range p.Sessions {
		n += s.Size
	}
	return n
}

// List returns every project under root, each with its sessions populated.
// Projects are sorted by most-recent activity first.
func List(root string) ([]Project, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read projects root %q: %w", root, err)
	}
	var out []Project
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		p := Project{
			EncodedName: e.Name(),
			Path:        filepath.Join(root, e.Name()),
		}
		sessions, err := listSessions(p.Path)
		if err != nil {
			return nil, err
		}
		p.Sessions = sessions
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].LastActivity().After(out[j].LastActivity())
	})
	return out, nil
}

// Find resolves a user-supplied project identifier (either an encoded
// directory name or a decoded path) against the projects under root.
func Find(root, ident string) (Project, error) {
	all, err := List(root)
	if err != nil {
		return Project{}, err
	}
	want := strings.ReplaceAll(ident, "/", "-")
	want = strings.TrimPrefix(want, "-")
	for _, p := range all {
		name := strings.TrimPrefix(p.EncodedName, "-")
		if name == want || p.EncodedName == ident {
			return p, nil
		}
	}
	return Project{}, fmt.Errorf("no project matches %q under %s", ident, root)
}

func listSessions(dir string) ([]Session, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read project dir %q: %w", dir, err)
	}
	var out []Session
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return nil, err
		}
		out = append(out, Session{
			ID:       strings.TrimSuffix(e.Name(), ".jsonl"),
			Path:     filepath.Join(dir, e.Name()),
			Size:     info.Size(),
			Modified: info.ModTime(),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Modified.After(out[j].Modified)
	})
	return out, nil
}
