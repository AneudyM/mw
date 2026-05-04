// Package stats aggregates information across many projects.
package stats

import (
	"time"

	"github.com/aneudym/mw/claude-projects-monitor/internal/projects"
)

// Summary is the top-level aggregate.
type Summary struct {
	Projects     int
	Sessions     int
	TotalBytes   int64
	LastActivity time.Time
	Top          []ProjectStat // largest projects by session count
}

// ProjectStat is the per-project entry inside Summary.Top.
type ProjectStat struct {
	Name     string
	Sessions int
	Bytes    int64
	Last     time.Time
}

// Aggregate computes a Summary for the given projects. The Top slice is
// limited to topN entries (or all of them, if topN <= 0).
func Aggregate(ps []projects.Project, topN int) Summary {
	s := Summary{Projects: len(ps)}
	for _, p := range ps {
		s.Sessions += len(p.Sessions)
		s.TotalBytes += p.TotalSize()
		if last := p.LastActivity(); last.After(s.LastActivity) {
			s.LastActivity = last
		}
		s.Top = append(s.Top, ProjectStat{
			Name:     p.EncodedName,
			Sessions: len(p.Sessions),
			Bytes:    p.TotalSize(),
			Last:     p.LastActivity(),
		})
	}
	sortByCountDesc(s.Top)
	if topN > 0 && len(s.Top) > topN {
		s.Top = s.Top[:topN]
	}
	return s
}

func sortByCountDesc(xs []ProjectStat) {
	for i := 1; i < len(xs); i++ {
		for j := i; j > 0 && xs[j-1].Sessions < xs[j].Sessions; j-- {
			xs[j-1], xs[j] = xs[j], xs[j-1]
		}
	}
}
