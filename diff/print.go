package diff

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// PrintMultiFileDiff prints a multi-file diff in unified diff format.
//
// To avoid holding the whole diff in memory, use WriteMultiFileDiff.
func PrintMultiFileDiff(ds []*FileDiff) ([]byte, error) {
	return printExact(func(w io.Writer) error { return WriteMultiFileDiff(w, ds) })
}

// PrintFileDiff prints a FileDiff in unified diff format.
//
// To avoid holding the whole diff in memory, use WriteFileDiff.
func PrintFileDiff(d *FileDiff) ([]byte, error) {
	return printExact(func(w io.Writer) error { return WriteFileDiff(w, d) })
}

// PrintHunks prints diff hunks in unified diff format.
func PrintHunks(hunks []*Hunk) ([]byte, error) {
	return printExact(func(w io.Writer) error { return (&printer{w: w}).hunks(hunks) })
}

// printExact calls write twice: first to measure the output, then to fill a
// buffer of exactly that size. This keeps the peak memory of printing a large
// diff at the size of the output, instead of the up to 2x that a growing buffer
// needs while it copies itself into a larger allocation.
func printExact(write func(io.Writer) error) ([]byte, error) {
	var n countingWriter
	if err := write(&n); err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	buf := bytes.NewBuffer(make([]byte, 0, n))
	if err := write(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// countingWriter counts the bytes written to it and discards them.
type countingWriter int

func (c *countingWriter) Write(p []byte) (int, error) {
	*c += countingWriter(len(p))
	return len(p), nil
}

// WriteMultiFileDiff writes a multi-file diff in unified diff format to w.
//
// Hunk bodies are written to w directly, without intermediate copies, so many
// small writes go to w. Wrap w in a bufio.Writer if those writes are costly.
func WriteMultiFileDiff(w io.Writer, ds []*FileDiff) error {
	p := printer{w: w}
	for _, d := range ds {
		if err := p.fileDiff(d); err != nil {
			return err
		}
	}
	return nil
}

// WriteFileDiff writes a FileDiff in unified diff format to w.
//
// Hunk bodies are written to w directly, without intermediate copies, so many
// small writes go to w. Wrap w in a bufio.Writer if those writes are costly.
//
// File names that contain control characters are written as quoted strings.
func WriteFileDiff(w io.Writer, d *FileDiff) error {
	return (&printer{w: w}).fileDiff(d)
}

// printer writes diffs to w. Header lines are assembled in line, which is
// reused across lines so that printing does not allocate per line.
type printer struct {
	w    io.Writer
	line []byte
}

func (p *printer) fileDiff(d *FileDiff) error {
	for _, xheader := range d.Extended {
		p.line = append(append(p.line[:0], xheader...), '\n')
		if err := p.writeLine(); err != nil {
			return err
		}
	}

	// FileDiff is added/deleted file
	// No further hunks printing needed
	if d.NewName == "" {
		_, err := fmt.Fprintf(p.w, onlyInMessage, filepath.Dir(d.OrigName), filepath.Base(d.OrigName))
		return err
	}

	if d.Hunks == nil {
		return nil
	}

	if err := p.fileHeader("--- ", d.OrigName, d.OrigTime); err != nil {
		return err
	}
	if err := p.fileHeader("+++ ", d.NewName, d.NewTime); err != nil {
		return err
	}
	return p.hunks(d.Hunks)
}

func (p *printer) fileHeader(prefix string, filename string, timestamp *time.Time) error {
	p.line = append(p.line[:0], prefix...)
	if needsQuoting(filename) {
		p.line = appendQuotedFilename(p.line, filename)
	} else {
		p.line = append(p.line, filename...)
	}
	if timestamp != nil {
		p.line = timestamp.AppendFormat(append(p.line, '\t'), diffTimeFormatLayout)
	}
	p.line = append(p.line, '\n')
	return p.writeLine()
}

// needsQuoting reports whether filename cannot be written verbatim in a
// "---" or "+++" header: control characters would end the name early (a tab
// starts the timestamp) or split the header, and a leading quote would make
// the name look quoted when it is read back.
func needsQuoting(filename string) bool {
	if strings.HasPrefix(filename, `"`) {
		return true
	}
	for i := 0; i < len(filename); i++ {
		if c := filename[i]; c < ' ' || c == 0x7f {
			return true
		}
	}
	return false
}

// appendQuotedFilename appends filename to dst as a C-style quoted string, the
// way git quotes unusual paths, so that the parser unquotes it back to
// filename.
func appendQuotedFilename(dst []byte, filename string) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(filename); {
		c := filename[i]
		switch c {
		case '"', '\\':
			dst = append(dst, '\\', c)
		case '\a':
			dst = append(dst, '\\', 'a')
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		case '\v':
			dst = append(dst, '\\', 'v')
		default:
			if c < ' ' || c == 0x7f {
				dst = appendOctal(dst, c)
			} else if r, size := utf8.DecodeRuneInString(filename[i:]); r == utf8.RuneError && size == 1 {
				// Not valid UTF-8: escape the byte so it survives unquoting.
				dst = appendOctal(dst, c)
			} else {
				dst = append(dst, filename[i:i+size]...)
				i += size
				continue
			}
		}
		i++
	}
	return append(dst, '"')
}

// appendOctal appends c as a three-digit octal escape such as \033.
func appendOctal(dst []byte, c byte) []byte {
	return append(dst, '\\', '0'+c>>6, '0'+(c>>3)&7, '0'+c&7)
}

func (p *printer) hunks(hunks []*Hunk) error {
	for _, hunk := range hunks {
		// Equivalent to fmt.Sprintf(hunkHeader, ...) plus the section.
		p.line = append(p.line[:0], "@@ -"...)
		p.line = strconv.AppendInt(p.line, int64(hunk.OrigStartLine), 10)
		p.line = append(p.line, ',')
		p.line = strconv.AppendInt(p.line, int64(hunk.OrigLines), 10)
		p.line = append(p.line, " +"...)
		p.line = strconv.AppendInt(p.line, int64(hunk.NewStartLine), 10)
		p.line = append(p.line, ',')
		p.line = strconv.AppendInt(p.line, int64(hunk.NewLines), 10)
		p.line = append(p.line, " @@"...)
		if hunk.Section != "" {
			p.line = append(append(p.line, ' '), hunk.Section...)
		}
		p.line = append(p.line, '\n')
		if err := p.writeLine(); err != nil {
			return err
		}

		if hunk.OrigNoNewlineAt == 0 {
			if _, err := p.w.Write(hunk.Body); err != nil {
				return err
			}
		} else {
			if _, err := p.w.Write(hunk.Body[:hunk.OrigNoNewlineAt]); err != nil {
				return err
			}
			if err := p.noNewlineMessage(); err != nil {
				return err
			}
			if _, err := p.w.Write(hunk.Body[hunk.OrigNoNewlineAt:]); err != nil {
				return err
			}
		}

		if !bytes.HasSuffix(hunk.Body, []byte{'\n'}) {
			p.line = append(append(append(p.line[:0], '\n'), noNewlineMessage...), '\n')
			if err := p.writeLine(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *printer) noNewlineMessage() error {
	p.line = append(append(p.line[:0], noNewlineMessage...), '\n')
	return p.writeLine()
}

func (p *printer) writeLine() error {
	_, err := p.w.Write(p.line)
	return err
}
