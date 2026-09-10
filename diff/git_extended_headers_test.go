package diff

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestParseFileDiffPreservesCRLFGitExtendedHeaders(t *testing.T) {
	input := []byte("diff --git a/old name b/new name\r\nsimilarity index 100%\r\nrename from old name\r\nrename to new name\r\n")
	wantExtended := []string{
		"diff --git a/old name b/new name\r",
		"similarity index 100%\r",
		"rename from old name\r",
		"rename to new name\r",
	}

	fd, err := ParseFileDiffOptions(input, ParseOptions{KeepCR: true})
	if err != nil {
		t.Fatal(err)
	}
	// Keep the existing edge behavior: with KeepCR, reconstruction
	// finds the first ambiguous path but not the second.
	if fd.OrigName != "a/old name" || fd.NewName != "" {
		t.Errorf("got names %q -> %q, want %q -> empty", fd.OrigName, fd.NewName, "a/old name")
	}
	if !reflect.DeepEqual(fd.Extended, wantExtended) {
		t.Errorf("extended headers changed:\nwant: %q\n got: %q", wantExtended, fd.Extended)
	}
}

func TestHandleEmptyRejectsMalformedGitExtendedHeaders(t *testing.T) {
	tests := []struct {
		name     string
		extended []string
	}{
		{name: "empty"},
		{name: "non-git", extended: []string{"diff --ruN a/f b/f", "old mode 100644", "new mode 100755"}},
		{name: "reversed mode pair", extended: []string{"diff --git a/f b/f", "new mode 100755", "old mode 100644"}},
		{name: "copy prefixes missing spaces", extended: []string{"diff --git a/f b/g", "similarity index 100%", "copy from", "copy to"}},
		{name: "unknown extra header", extended: []string{"diff --git a/f b/f", "old mode 100644", "new mode 100755", "x-header value"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fd := &FileDiff{Extended: test.extended}
			if handleEmpty(fd) {
				t.Fatal("handleEmpty accepted malformed headers")
			}
			if fd.OrigName != "" || fd.NewName != "" {
				t.Errorf("handleEmpty changed names to %q -> %q", fd.OrigName, fd.NewName)
			}
		})
	}
}

func TestMalformedGitExtendedHeadersStillReturnOriginalParseError(t *testing.T) {
	input := []byte("diff --git a/f b/f\nold mode 100644\nunknown header\n")
	r := NewFileDiffReader(bytes.NewReader(input))
	fd, err := r.ReadAllHeaders()
	var parseErr *ParseError
	if !errors.As(err, &parseErr) || parseErr.Err != ErrExtendedHeadersEOF {
		t.Fatalf("got error %v, want ErrExtendedHeadersEOF", err)
	}
	wantExtended := []string{"diff --git a/f b/f", "old mode 100644", "unknown header"}
	if !reflect.DeepEqual(fd.Extended, wantExtended) {
		t.Errorf("extended headers changed:\nwant: %q\n got: %q", wantExtended, fd.Extended)
	}
}
