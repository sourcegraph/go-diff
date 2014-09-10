package diff

import "time"

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

// A Hunk represents a series of changes (additions or deletions) in a
// file's unified diff.
type Hunk struct {
	OrigStartLine int // starting line number in original file
	OrigLines     int // number of lines the hunk applies to in the original file

	NewStartLine int // starting line number in new file
	NewLines     int // number of lines the hunk applies to in the new file

	Section string // optional section heading

	Body []byte // hunk body (lines prefixed with '-', '+', or ' ')
}

var (
	hunkPrefix = []byte("@@ ")
)

const hunkHeader = "@@ -%d,%d +%d,%d @@"

// diffTimeFormat is the time format string for unified diff file
// header timestamps. See
// http://www.gnu.org/software/diffutils/manual/html_node/Detailed-Unified.html.
const diffTimeFormat = "2006-01-02 15:04:05.000000000 -0700"
