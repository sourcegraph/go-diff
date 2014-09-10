package diff

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

// ParseFileDiff parses a file unified diff.
func ParseFileDiff(diff []byte) (*FileDiff, error) {
	return NewFileDiffReader(bytes.NewReader(diff)).Read()
}

// NewFileDiffReader returns a new FileDiffReader that reads a file
// unified diff.
func NewFileDiffReader(r io.Reader) *FileDiffReader {
	return &FileDiffReader{scanner: bufio.NewScanner(r)}
}

// FileDiffReader reads a unified file diff.
type FileDiffReader struct {
	line    int
	offset  int64
	scanner *bufio.Scanner

	// fileHeaderLine is the first file header line, set by
	// ReadExtendedHeaders if it encroaches on a file header line
	// (which it must to detect when extended headers are done).
	fileHeaderLine []byte
}

// Read reads a file unified diff, including headers and hunks, from r.
func (r *FileDiffReader) Read() (*FileDiff, error) {
	d, err := r.ReadAllHeaders()
	if err != nil {
		return nil, err
	}

	d.Hunks, err = r.HunksReader().ReadAllHunks()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// ReadAllHeaders reads the file headers and extended headers (if any)
// from a file unified diff. It does not read hunks, and the returned
// FileDiff's Hunks field is nil. To read the hunks, call the
// (*FileDiffReader).HunksReader() method to get a HunksReader and
// read hunks from that.
func (r *FileDiffReader) ReadAllHeaders() (*FileDiff, error) {
	var err error
	fd := &FileDiff{}

	fd.Extended, err = r.ReadExtendedHeaders()
	if err != nil {
		return nil, err
	}

	fd.OrigName, fd.NewName, fd.OrigTime, fd.NewTime, err = r.ReadFileHeaders()
	if err != nil {
		return nil, err
	}

	return fd, nil
}

// HunksReader returns a new HunksReader that reads hunks from r. The
// HunksReader's line and offset (used in error messages) is set to
// start where the file diff header ended (which means errors have the
// correct position information).
func (r *FileDiffReader) HunksReader() *HunksReader {
	return &HunksReader{
		line:    r.line,
		offset:  r.offset,
		scanner: r.scanner,
	}
}

// ReadFileHeaders reads the unified file diff header (the lines that
// start with "---" and "+++" with the orig/new file names and
// timestamps).
func (r *FileDiffReader) ReadFileHeaders() (origName, newName string, origTimestamp, newTimestamp *time.Time, err error) {
	origName, origTimestamp, err = r.readOneFileHeader([]byte("--- "))
	if err != nil {
		return "", "", nil, nil, err
	}

	newName, newTimestamp, err = r.readOneFileHeader([]byte("+++ "))
	if err != nil {
		return "", "", nil, nil, err
	}

	return origName, newName, origTimestamp, newTimestamp, nil
}

// readOneFileHeader reads one of the file headers (prefix should be
// either "+++ " or "--- ").
func (r *FileDiffReader) readOneFileHeader(prefix []byte) (filename string, timestamp *time.Time, err error) {
	var line []byte

	if r.fileHeaderLine == nil {
		ok := r.scanner.Scan()
		if !ok {
			return "", nil, &ParseError{r.line, r.offset, ErrNoFileHeader}
		}
		line = r.scanner.Bytes()
	} else {
		line = r.fileHeaderLine
		r.fileHeaderLine = nil
	}

	if !bytes.HasPrefix(line, prefix) {
		return "", nil, &ParseError{r.line, r.offset, ErrBadFileHeader}
	}

	r.offset += int64(len(line))
	r.line++
	line = line[len(prefix):]

	parts := bytes.SplitN(line, []byte("\t"), 2)
	filename = string(parts[0])
	if len(parts) == 2 {
		// Timestamp is optional, but this header has it.
		ts, err := time.Parse(diffTimeFormat, string(parts[1]))
		if err != nil {
			return "", nil, err
		}
		timestamp = &ts
	}

	return filename, timestamp, err
}

// ReadExtendedHeaders reads the extended header lines, if any, from a
// unified diff file (e.g., git's "diff --git a/foo.go b/foo.go", "new
// mode <mode>", "rename from <path>", etc.).
func (r *FileDiffReader) ReadExtendedHeaders() ([]string, error) {
	var xheaders []string
	for {
		ok := r.scanner.Scan()
		if !ok {
			return nil, &ParseError{r.line, r.offset, ErrExtendedHeadersEOF}
		}

		line := r.scanner.Bytes()

		if bytes.HasPrefix(line, []byte("--- ")) {
			// We've reached the file header.
			r.fileHeaderLine = line // pass to readOneFileHeader (see fileHeaderLine field doc)
			return xheaders, nil
		}

		r.line++
		r.offset += int64(len(line))
		xheaders = append(xheaders, string(line))
	}
}

var (
	// ErrNoFileHeader is when a file unified diff has no file header
	// (i.e., the lines that begin with "---" and "+++").
	ErrNoFileHeader = errors.New("expected file header, got EOF")

	// ErrBadFileHeader is when a file unified diff has a malformed
	// file header (i.e., the lines that begin with "---" and "+++").
	ErrBadFileHeader = errors.New("bad file header")

	// ErrExtendedHeadersEOF is when an EOF was encountered while reading extended file headers, which means that there were no ---/+++ headers encountered before hunks (if any) began.
	ErrExtendedHeadersEOF = errors.New("expected file header while reading extended headers, got EOF")
)

// ParseHunks parses hunks from a unified diff. The diff must consist
// only of hunks and not include a file header; if it has a file
// header, use ParseFileDiff.
func ParseHunks(diff []byte) ([]*Hunk, error) {
	r := NewHunksReader(bytes.NewReader(diff))
	hunks, err := r.ReadAllHunks()
	if err != nil {
		return nil, err
	}
	return hunks, nil
}

// NewHunksReader returns a new HunksReader that reads unified diff hunks
// from r.
func NewHunksReader(r io.Reader) *HunksReader {
	return &HunksReader{scanner: bufio.NewScanner(r)}
}

// A HunksReader reads hunks from a unified diff.
type HunksReader struct {
	line    int
	offset  int64
	hunk    *Hunk
	scanner *bufio.Scanner

	nextHunkHeaderLine []byte
}

// ReadHunk reads one hunk from r.
func (r *HunksReader) ReadHunk() (*Hunk, error) {
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
func (r *HunksReader) ReadAllHunks() ([]*Hunk, error) {
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
