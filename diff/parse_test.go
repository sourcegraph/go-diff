package diff

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadQuotedFilename_Success(t *testing.T) {
	tests := []struct {
		input, value, remainder string
	}{
		{input: `""`, value: "", remainder: ""},
		{input: `"aaa"`, value: "aaa", remainder: ""},
		{input: `"aaa" bbb`, value: "aaa", remainder: " bbb"},
		{input: `"aaa" "bbb" ccc`, value: "aaa", remainder: ` "bbb" ccc`},
		{input: `"\""`, value: "\"", remainder: ""},
		{input: `"uh \"oh\""`, value: "uh \"oh\"", remainder: ""},
		{input: `"uh \\"oh\\""`, value: "uh \\", remainder: `oh\\""`},
		{input: `"uh \\\"oh\\\""`, value: "uh \\\"oh\\\"", remainder: ""},
	}
	for _, tc := range tests {
		value, remainder, err := readQuotedFilename(tc.input)
		if err != nil {
			t.Errorf("readQuotedFilename(`%s`): expected success, got '%s'", tc.input, err)
		} else if value != tc.value || remainder != tc.remainder {
			t.Errorf("readQuotedFilename(`%s`): expected `%s` and `%s`, got `%s` and `%s`", tc.input, tc.value, tc.remainder, value, remainder)
		}
	}
}

func TestReadQuotedFilename_Error(t *testing.T) {
	tests := []string{
		// Doesn't start with a quote
		``,
		`foo`,
		` "foo"`,
		// Missing end quote
		`"`,
		`"\"`,
		// "\x" is not a valid Go string literal escape
		`"\xxx"`,
	}
	for _, input := range tests {
		_, _, err := readQuotedFilename(input)
		if err == nil {
			t.Errorf("readQuotedFilename(`%s`): expected error", input)
		}
	}
}

func TestParseDiffGitArgs_Success(t *testing.T) {
	tests := []struct {
		input, first, second string
	}{
		{input: `aaa bbb`, first: "aaa", second: "bbb"},
		{input: `"aaa" bbb`, first: "aaa", second: "bbb"},
		{input: `aaa "bbb"`, first: "aaa", second: "bbb"},
		{input: `"aaa" "bbb"`, first: "aaa", second: "bbb"},
		{input: `1/a 2/z`, first: "1/a", second: "2/z"},
		{input: `1/hello world 2/hello world`, first: "1/hello world", second: "2/hello world"},
		{input: `"new\nline" and spaces`, first: "new\nline", second: "and spaces"},
		{input: `a/existing file with spaces "b/new, complicated\nfilen\303\270me"`, first: "a/existing file with spaces", second: "b/new, complicated\nfilen\303\270me"},
	}
	for _, tc := range tests {
		first, second, success := parseDiffGitArgs(tc.input)
		if !success {
			t.Errorf("`diff --git %s`: expected success", tc.input)
		} else if first != tc.first || second != tc.second {
			t.Errorf("`diff --git %s`: expected `%s` and `%s`, got `%s` and `%s`", tc.input, tc.first, tc.second, first, second)
		}
	}
}

func TestParseDiffGitArgs_Unsuccessful(t *testing.T) {
	tests := []string{
		``,
		`hello_world.txt`,
		`word `,
		` word`,
		`"a/bad_quoting b/bad_quoting`,
		`a/bad_quoting "b/bad_quoting`,
		`a/bad_quoting b/bad_quoting"`,
		`"a/bad_quoting b/bad_quoting"`,
		`"a/bad""b/bad"`,
		`"a/bad" "b/bad" "c/bad"`,
		`a/bad "b/bad" "c/bad"`,
	}
	for _, input := range tests {
		first, second, success := parseDiffGitArgs(input)
		if success {
			t.Errorf("`diff --git %s`: expected unsuccessful; got `%s` and `%s`", input, first, second)
		}
	}
}

// virtualDiff implements io.Reader to generate a 'virtual' diff file.
type virtualDiff struct {
	// diffFileHeader contains diffFile and hunk header
	diffFileHeader string
	// diffFileRepeats is number of times to repeat diffFileHeader
	diffFileRepeats int

	// line is a hunk content/body line
	line string
	// lineRepeats is number of times to repeat line
	lineRepeats int

	currentFileHeader       int
	currentFileHeaderOffset int

	currentLine       int
	currentLineOffset int
}

func (vd *virtualDiff) reset() {
	vd.currentFileHeader = 0
	vd.currentFileHeaderOffset = 0

	vd.currentLine = 0
	vd.currentLineOffset = 0
}

func (vd *virtualDiff) Read(p []byte) (n int, err error) {

	for {
		if vd.currentFileHeader >= vd.diffFileRepeats {
			return 0, io.EOF
		}

		// read from file header
		if vd.currentFileHeaderOffset < len(vd.diffFileHeader) {
			n = copy(p, vd.diffFileHeader[vd.currentFileHeaderOffset:])
			vd.currentFileHeaderOffset += n
			return
		}

		// read from line
		if vd.currentLineOffset < len(vd.line) {
			n = copy(p, vd.line[vd.currentLineOffset:])
			vd.currentLineOffset += n
			return
		}

		// must be at end-of-line
		vd.currentLine++
		vd.currentLineOffset = 0

		// must be at end-of-diff-file
		if vd.currentLine >= vd.lineRepeats {
			vd.currentFileHeader++
			vd.currentFileHeaderOffset = 0
			vd.currentLine = 0
			vd.currentLineOffset = 0
		}
	}
}

// emptyContentHandler implements ContentHandler and does nothing
type emptyContentHandler struct {
	fileNo int
}

func (h *emptyContentHandler) StartFile(fileDiff *FileDiff) error {
	fmt.Println(fmt.Sprintf("StartFile: %d", h.fileNo))
	h.fileNo++
	return nil
}
func (h *emptyContentHandler) EndFile(fileDiff *FileDiff) error {
	return nil
}
func (h *emptyContentHandler) StartHunk(hunk *Hunk) error {
	return nil
}
func (h *emptyContentHandler) EndHunk(hunk *Hunk) error {
	return nil
}
func (h *emptyContentHandler) HunkLine(hunk *Hunk, line []byte, eol bool) error {
	return nil
}

func BenchmarkReadDiffWithContentHandler_InfiniteFile(b *testing.B) {
	a := assert.New(b)

	virtualDiff := &virtualDiff{}
	virtualDiff.diffFileHeader = `diff --git a/README.md b/README.md
index aa4de15..7c048ab 100644
--- oldname	2009-10-11 15:12:20.000000000 +0000
+++ newname	2009-10-11 15:12:30.000000000 +0000
@@ -1,1 +1,1 @@
`

	// each file is 100GB
	virtualDiff.line = `+` + strings.Repeat("a", 1000000) + "\n" // 1MB
	virtualDiff.lineRepeats = 100000
	// total diff size is 1TB
	virtualDiff.diffFileRepeats = 10

	// total virtual diff size is 1TB
	// as we should read in streaming mode, no memory should be used
	err := ReadDiffWithContentHandler(virtualDiff, &emptyContentHandler{})
	a.NoError(err)

	// in contrast, try running without content handler and see GBs of memory used :)
	// err := ReadDiffWithContentHandler(virtualDiff, nil)
	// a.NoError(err)
}

func TestReadDiffWithContentHandler_LongHunkLines(t *testing.T) {
	a := assert.New(t)

	virtualDiff := &virtualDiff{}
	virtualDiff.diffFileHeader = `diff --git a/README.md b/README.md
index aa4de15..7c048ab 100644
--- oldname	2009-10-11 15:12:20.000000000 +0000
+++ newname	2009-10-11 15:12:30.000000000 +0000
@@ -1,1 +1,1 @@
`

	virtualDiff.line = `+` + strings.Repeat("a", 1000000) + "\n" // 1MB
	virtualDiff.lineRepeats = 1
	virtualDiff.diffFileRepeats = 1

	// first, run without contentHandler
	fileDiffs, err := NewMultiFileDiffReader(virtualDiff).ReadAllFiles()
	a.NoError(err)
	a.Equal(1, len(fileDiffs))
	a.Equal(1, len(fileDiffs[0].Hunks))
	a.Equal([]byte(virtualDiff.line), fileDiffs[0].Hunks[0].Body)

	// second, read with contentHandler
	virtualDiff.reset()
	h := &domContentHandler{}
	err = ReadDiffWithContentHandler(virtualDiff, h)
	a.NoError(err)
	a.Equal(1, len(h.fileDiffs))
	a.Equal(1, len(h.fileDiffs[0].Hunks))
	a.Equal([]byte(virtualDiff.line), h.fileDiffs[0].Hunks[0].Body)
}
