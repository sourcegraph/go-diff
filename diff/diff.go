package diff

import (
	"bytes"
	"time"
)

// A FileDiff represents a unified diff for a single file.
//
// A file unified diff has a header that resembles the following:
//
//  --- oldname	2009-10-11 15:12:20.000000000 -0700
//  +++ newname	2009-10-11 15:12:30.000000000 -0700
type FileDiff struct {
	OrigName string     // the original name of the file
	OrigTime *time.Time // the original timestamp (nil if not present)

	NewName string     // the new name of the file (often same as OrigName)
	NewTime *time.Time // the new timestamp (nil if not present)

	Extended []string // extended header lines (e.g., git's "new mode <mode>", "rename from <path>", etc.)

	Hunks []*Hunk // hunks that were changed from orig to new
}

// Stat computes the number of lines added/changed/deleted in all
// hunks in this file's diff.
func (d *FileDiff) Stat() Stat {
	total := Stat{}
	for _, h := range d.Hunks {
		total.add(h.Stat())
	}
	return total
}

// A Hunk represents a series of changes (additions or deletions) in a
// file's unified diff.
type Hunk struct {
	OrigStartLine   int // starting line number in original file
	OrigLines       int // number of lines the hunk applies to in the original file
	OrigNoNewlineAt int // if > 0, then the original file had a 'No newline at end of file' mark at this offset

	NewStartLine int // starting line number in new file
	NewLines     int // number of lines the hunk applies to in the new file

	Section string // optional section heading

	Body []byte // hunk body (lines prefixed with '-', '+', or ' ')
}

// Stat computes the number of lines added/changed/deleted in this
// hunk.
func (h *Hunk) Stat() Stat {
	lines := bytes.Split(h.Body, []byte{'\n'})
	var last byte
	st := Stat{}
	for _, line := range lines {
		if len(line) == 0 {
			last = 0
			continue
		}
		switch line[0] {
		case '-':
			if last == '+' {
				st.Added--
				st.Changed++
				last = 0 // next line can't change this one since this is already a change
			} else {
				st.Deleted++
				last = line[0]
			}
		case '+':
			if last == '-' {
				st.Deleted--
				st.Changed++
				last = 0 // next line can't change this one since this is already a change
			} else {
				st.Added++
				last = line[0]
			}
		default:
			last = 0
		}
	}
	return st
}

var (
	hunkPrefix = []byte("@@ ")
)

const hunkHeader = "@@ -%d,%d +%d,%d @@"

// diffTimeFormat is the time format string for unified diff file
// header timestamps. See
// http://www.gnu.org/software/diffutils/manual/html_node/Detailed-Unified.html.
const diffTimeFormat = "2006-01-02 15:04:05.000000000 -0700"

// A Stat is a diff stat that represents the number of lines
// added/changed/deleted.
type Stat struct {
	Added, Changed, Deleted int // numbers of lines
}

func (s *Stat) add(o Stat) {
	s.Added += o.Added
	s.Changed += o.Changed
	s.Deleted += o.Deleted
}
