package diff

import (
	"bytes"
	"fmt"
)

// Print prints a Diff in unified diff format.
func Print(diff *Diff) ([]byte, error) {
	var buf bytes.Buffer
	for _, hunk := range diff.Hunks {
		_, err := fmt.Fprintf(&buf,
			"@@ -%d,%d +%d,%d @@", hunk.OrigStartLine, hunk.OrigLines, hunk.NewStartLine, hunk.NewLines,
		)
		if err != nil {
			return nil, err
		}
		if hunk.Section != "" {
			_, err := fmt.Fprint(&buf, " ", hunk.Section)
			if err != nil {
				return nil, err
			}
		}
		if _, err := fmt.Fprintln(&buf); err != nil {
			return nil, err
		}
		if _, err := buf.Write(hunk.Body); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
