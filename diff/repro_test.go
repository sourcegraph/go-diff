package diff

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestHandleEmpty_KeepCR(t *testing.T) {
	// A diff with CRLF line endings
	input := "diff --git a/file b/file\r\n" +
		"new file mode 100644\r\n" +
		"index 0000000..e69de29\r\n"

	opts := ParseOptions{KeepCR: true}
	fds, err := ParseMultiFileDiffOptions([]byte(input), opts)
	if err != nil {
		t.Fatal(err)
	}

	if len(fds) != 1 {
		t.Fatalf("expected 1 FileDiff, got %d", len(fds))
	}

	fd := fds[0]
	// The problem: if KeepCR is true, NewName might be "b/file\r"
	wantNewName := "b/file"
	if fd.NewName != wantNewName {
		t.Errorf("expected NewName %q, got %q", wantNewName, fd.NewName)
	}

	// Also check Extended headers
	wantExtended := []string{
		"diff --git a/file b/file\r",
		"new file mode 100644\r",
		"index 0000000..e69de29\r",
	}
	if !reflect.DeepEqual(fd.Extended, wantExtended) {
		t.Errorf("Extended headers mismatch.\nwant: %q\ngot:  %q", wantExtended, fd.Extended)
	}
}
func TestParseOnlyIn_KeepCR(t *testing.T) {
	// A diff with an "Only in" message and CRLF line ending
	input := "Only in /tmp: file\r\n"

	opts := ParseOptions{KeepCR: true}
	fds, err := ParseMultiFileDiffOptions([]byte(input), opts)
	if err != nil {
		t.Fatal(err)
	}

	if len(fds) != 1 {
		t.Fatalf("expected 1 FileDiff, got %d", len(fds))
	}

	fd := fds[0]
	// The path should be /tmp/file (or \tmp\file on windows, but filepath.Join is used)
	// Actually, the code uses filepath.Join(string(source), string(filename))
	wantOrigName := filepath.Join("/tmp", "file")
	if fd.OrigName != wantOrigName {
		t.Errorf("expected OrigName %q, got %q", wantOrigName, fd.OrigName)
	}

	if fd.NewName != "" {
		t.Errorf("expected empty NewName, got %q", fd.NewName)
	}
}

func TestNormalizeHeader_KeepCR(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantH   string
		wantS   string
		wantErr bool
	}{
		{
			name:  "NoSection_CR",
			input: "@@ -1,1 +1,1 @@\r",
			wantH: "@@ -1,1 +1,1 @@",
			wantS: "",
		},
		{
			name:  "WithSection_CR",
			input: "@@ -1,1 +1,1 @@ some section\r",
			wantH: "@@ -1,1 +1,1 @@ some section",
			wantS: "some section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotH, gotS, err := normalizeHeader(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotH != tt.wantH {
				t.Errorf("normalizeHeader() gotH = %q, want %q", gotH, tt.wantH)
			}
			if gotS != tt.wantS {
				t.Errorf("normalizeHeader() gotS = %q, want %q", gotS, tt.wantS)
			}
		})
	}
}

func TestNoNewline_KeepCR(t *testing.T) {
	// A diff with CRLF line endings and a "No newline at end of file" marker
	input := "--- a/file\r\n" +
		"+++ b/file\r\n" +
		"@@ -1 +1 @@\r\n" +
		"-a\r\n" +
		"\\ No newline at end of file\r\n" +
		"+b\r\n"

	opts := ParseOptions{KeepCR: true}
	fds, err := ParseMultiFileDiffOptions([]byte(input), opts)
	if err != nil {
		t.Fatal(err)
	}

	if len(fds) != 1 {
		t.Fatalf("expected 1 FileDiff, got %d", len(fds))
	}

	fd := fds[0]
	if len(fd.Hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(fd.Hunks))
	}

	h := fd.Hunks[0]
	// The body should contain the lines exactly as they are (including \r\n)
	// but the "No newline" marker should NOT be in the body.
	wantBody := "-a\r\n+b\r\n"
	if string(h.Body) != wantBody {
		t.Errorf("expected body %q, got %q", wantBody, string(h.Body))
	}

	if h.OrigNoNewlineAt == 0 {
		t.Error("expected OrigNoNewlineAt to be set")
	}
}
