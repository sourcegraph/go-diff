package diff

// TODO(sqs): support parsing diffs that have hunks from multiple files (using the --- / +++ header syntax)

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Parse parses a unified diff.
func Parse(diff []byte) (*Diff, error) {
	r := NewReader(bytes.NewReader(diff))
	hunks, err := r.ReadAllHunks()
	if err != nil {
		return nil, err
	}
	return &Diff{Hunks: hunks}, nil
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{scanner: bufio.NewScanner(r)}
}

// A Reader reads hunks from a unified diff.
type Reader struct {
	line    int
	offset  int64
	hunk    *Hunk
	scanner *bufio.Scanner

	nextHunkHeaderLine []byte
}

// ReadHunk reads one hunk from r.
func (r *Reader) ReadHunk() (*Hunk, error) {
	r.hunk = nil
	var line []byte
	for {
		if r.nextHunkHeaderLine != nil {
			// Use stored hunk header line that was scanned in at the
			// completion of the previous hunk's ReadHunk.
			line = r.nextHunkHeaderLine
			r.nextHunkHeaderLine = nil
		} else {
			ok := r.scanner.Scan()
			if !ok {
				break
			}
			line = r.scanner.Bytes()
		}

		// Record position.
		r.line++
		r.offset += int64(len(line))

		if r.hunk == nil {
			// Check for presence of hunk header.
			if !bytes.HasPrefix(line, hunkPrefix) {
				return nil, &ParseError{r.line, r.offset, ErrNoHunkHeader}
			}

			// Parse hunk header.
			r.hunk = &Hunk{}
			items := []interface{}{
				&r.hunk.OrigStartLine, &r.hunk.OrigLines,
				&r.hunk.NewStartLine, &r.hunk.NewLines,
			}
			br := bytes.NewReader(line)
			n, err := fmt.Fscanf(br, hunkHeader, items...)
			if err != nil {
				return nil, err
			}
			if n < len(items) {
				return nil, &ParseError{r.line, r.offset, ErrBadHunkHeader}
			}

			// Any unread portion of the line is the (optional) section heading.
			if br.Len() > 0 {
				r.hunk.Section = string(bytes.TrimSpace(line[len(line)-br.Len():]))
			}
		} else {
			// Read hunk body line.
			if bytes.HasPrefix(line, hunkPrefix) {
				// Saw start of new hunk, so this hunk is
				// complete. But we've already read in the next hunk's
				// header, so we need to be sure that the next call to
				// ReadHunk starts with that header.
				r.nextHunkHeaderLine = line

				// Rewind position.
				r.line--
				r.offset -= int64(len(line))

				return r.hunk, nil
			}

			r.hunk.Body = append(r.hunk.Body, line...)
			r.hunk.Body = append(r.hunk.Body, '\n')
		}
	}
	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	// Final hunk is complete. But if we never saw a hunk in this call to ReadHunk, then it means we hit EOF.
	if r.hunk != nil {
		return r.hunk, nil
	}
	return nil, io.EOF
}

// ReadAllHunks reads all remaining hunks from r. A successful call
// returns err == nil, not err == EOF. Because ReadAllHunks is defined
// to read until EOF, it does not treat end of file as an error to be
// reported.
func (r *Reader) ReadAllHunks() ([]*Hunk, error) {
	var hunks []*Hunk
	for {
		hunk, err := r.ReadHunk()
		if err == io.EOF {
			return hunks, nil
		}
		if err != nil {
			return nil, err
		}
		hunks = append(hunks, hunk)
	}
}

// A ParseError is a description of a unified diff syntax error.
type ParseError struct {
	Line   int   // Line where the error occurred
	Offset int64 // Offset where the error occurred
	Err    error // The actual error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, char %d: %s", e.Line, e.Offset, e.Err)
}

var (
	// ErrNoHunkHeader indicates that a unified diff hunk header was
	// expected but not found during parsing.
	ErrNoHunkHeader = errors.New("no hunk header")

	// ErrBadHunkHeader indicates that a malformed unified diff hunk
	// header was encountered during parsing.
	ErrBadHunkHeader = errors.New("bad hunk header")
)
