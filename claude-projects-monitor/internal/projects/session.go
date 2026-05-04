package projects

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Entry is a minimal view of one line in a session transcript. Claude Code
// writes a richer JSON object per line; we only decode the fields we need
// and ignore everything else.
type Entry struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

// CountEntries returns the number of newline-delimited records in a session
// file. It is cheap relative to fully decoding each line.
func CountEntries(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open session %q: %w", path, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	n := 0
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		n++
	}
	if err := scanner.Err(); err != nil {
		return n, fmt.Errorf("scan session %q: %w", path, err)
	}
	return n, nil
}

// ReadEntries decodes every line of a session file into Entry values. Lines
// that fail to decode are skipped.
func ReadEntries(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open session %q: %w", path, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	var out []Entry
	for scanner.Scan() {
		var e Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		out = append(out, e)
	}
	if err := scanner.Err(); err != nil {
		return out, fmt.Errorf("scan session %q: %w", path, err)
	}
	return out, nil
}
