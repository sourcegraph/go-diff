package diff

// A Diff represents a unified diff.
type Diff struct {
	Hunks []*Hunk
}

// A Hunk represents a series of changes (additions or deletions) in a
// unified diff.
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
