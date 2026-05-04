// Package monitor polls the projects root and emits events when sessions
// appear, grow, or disappear. Polling avoids a third-party fsnotify
// dependency and is more than fast enough for a directory of transcripts.
package monitor

import (
	"context"
	"time"

	"github.com/aneudym/mw/claude-projects-monitor/internal/projects"
)

// EventKind classifies a change observed between two polls.
type EventKind int

const (
	SessionCreated EventKind = iota
	SessionUpdated
	SessionRemoved
)

func (k EventKind) String() string {
	switch k {
	case SessionCreated:
		return "created"
	case SessionUpdated:
		return "updated"
	case SessionRemoved:
		return "removed"
	}
	return "unknown"
}

// Event describes one observed change.
type Event struct {
	Kind    EventKind
	Project string
	Session string
	Path    string
	Size    int64
	When    time.Time
}

// Watch polls root every interval and writes Events to the returned channel
// until ctx is cancelled. The channel is closed when the loop exits.
func Watch(ctx context.Context, root string, interval time.Duration) (<-chan Event, <-chan error) {
	events := make(chan Event, 32)
	errs := make(chan error, 1)
	go func() {
		defer close(events)
		defer close(errs)
		prev, err := snapshot(root)
		if err != nil {
			errs <- err
			return
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				curr, err := snapshot(root)
				if err != nil {
					errs <- err
					return
				}
				diff(prev, curr, events)
				prev = curr
			}
		}
	}()
	return events, errs
}

type snap map[string]sessionState

type sessionState struct {
	project  string
	path     string
	size     int64
	modified time.Time
}

func snapshot(root string) (snap, error) {
	ps, err := projects.List(root)
	if err != nil {
		return nil, err
	}
	out := make(snap)
	for _, p := range ps {
		for _, s := range p.Sessions {
			out[s.Path] = sessionState{
				project:  p.EncodedName,
				path:     s.Path,
				size:     s.Size,
				modified: s.Modified,
			}
		}
	}
	return out, nil
}

func diff(prev, curr snap, out chan<- Event) {
	for path, c := range curr {
		p, ok := prev[path]
		if !ok {
			out <- Event{
				Kind:    SessionCreated,
				Project: c.project,
				Session: sessionID(path),
				Path:    path,
				Size:    c.size,
				When:    c.modified,
			}
			continue
		}
		if c.size != p.size || !c.modified.Equal(p.modified) {
			out <- Event{
				Kind:    SessionUpdated,
				Project: c.project,
				Session: sessionID(path),
				Path:    path,
				Size:    c.size,
				When:    c.modified,
			}
		}
	}
	for path, p := range prev {
		if _, ok := curr[path]; !ok {
			out <- Event{
				Kind:    SessionRemoved,
				Project: p.project,
				Session: sessionID(path),
				Path:    path,
				Size:    p.size,
				When:    time.Now(),
			}
		}
	}
}

func sessionID(path string) string {
	// trim ".jsonl" and any leading directory.
	name := path
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}
	if len(name) > 6 && name[len(name)-6:] == ".jsonl" {
		name = name[:len(name)-6]
	}
	return name
}
