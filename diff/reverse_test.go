package diff

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReverseHunks(t *testing.T) {
	tests := []struct {
		inputFile string
		wantFile  string
	}{
		{
			inputFile: "sample_hunks.diff",
			wantFile:  "sample_hunks.reversed",
		},
		{
			inputFile: "no_newline_new.diff",
			wantFile:  "no_newline_new.reversed",
		},
		{
			inputFile: "no_newline_orig.diff",
			wantFile:  "no_newline_orig.reversed",
		},
		{
			inputFile: "no_newline_both.diff",
			wantFile:  "no_newline_both.reversed",
		},
	}
	for _, test := range tests {
		inputData, err := os.ReadFile(filepath.Join("testdata", test.inputFile))
		if err != nil {
			t.Fatal(err)
		}
		wantData, err := os.ReadFile(filepath.Join("testdata", test.wantFile))
		if err != nil {
			t.Fatal(err)
		}
		input, err := ParseHunks(inputData)
		if err != nil {
			t.Fatal(err)
		}

		var reversed []*Hunk
		for _, in := range input {
			out, err := reverseHunk(in)
			if err != nil {
				// This should only fail if the Hunk data structure is inconsistent
				t.Errorf("%s: Unexpected reverseHunk() error: %s", test.inputFile, err)
			}
			reversed = append(reversed, out)
		}
		gotData, err := PrintHunks(reversed)
		if err != nil {
			t.Errorf("%s: PrintHunks of reversed data: %s", test.inputFile, err)
		}
		if !bytes.Equal(wantData, gotData) {
			t.Errorf("%s: Reversed hunk does not match expected.\nWant vs got:\n%s",
				test.inputFile, cmp.Diff(wantData, gotData))
		}
	}
}

func TestReverseFileDiff(t *testing.T) {
	tests := []struct {
		inputFile string
		wantFile  string
	}{
		{
			inputFile: "sample_file.diff",
			wantFile:  "sample_file.reversed",
		},
	}
	for _, test := range tests {
		inputData, err := os.ReadFile(filepath.Join("testdata", test.inputFile))
		if err != nil {
			t.Fatal(err)
		}
		wantData, err := os.ReadFile(filepath.Join("testdata", test.wantFile))
		if err != nil {
			t.Fatal(err)
		}
		input, err := ParseFileDiff(inputData)
		if err != nil {
			t.Fatal(err)
		}
		reversed, err := ReverseFileDiff(input)
		if err != nil {
			t.Errorf("%s: ReverseFileDiff: %s", test.inputFile, err)
		}
		gotData, err := PrintFileDiff(reversed)
		if err != nil {
			t.Errorf("%s: PrintFileDiff of reversed data: %s", test.inputFile, err)
		}
		if !bytes.Equal(wantData, gotData) {
			t.Errorf("%s: Reversed diff does not match expected.\nWant vs got:\n%s",
				test.inputFile, cmp.Diff(wantData, gotData))
		}
	}
}

func TestReverseMultiFileDiff(t *testing.T) {
	tests := []struct {
		inputFile string
		wantFile  string
	}{
		{
			inputFile: "sample_file.diff",
			wantFile:  "sample_file.reversed",
		},
		{
			inputFile: "sample_multi_file.diff",
			wantFile:  "sample_multi_file.reversed",
		},
	}
	for _, test := range tests {
		inputData, err := os.ReadFile(filepath.Join("testdata", test.inputFile))
		if err != nil {
			t.Fatal(err)
		}
		wantData, err := os.ReadFile(filepath.Join("testdata", test.wantFile))
		if err != nil {
			t.Fatal(err)
		}
		input, err := ParseMultiFileDiff(inputData)
		if err != nil {
			t.Fatal(err)
		}
		reversed, err := ReverseMultiFileDiff(input)
		if err != nil {
			t.Errorf("%s: ReverseMultiFileDiff: %s", test.inputFile, err)
		}
		gotData, err := PrintMultiFileDiff(reversed)
		if err != nil {
			t.Errorf("%s: PrintMultiFileDiff of reversed data: %s", test.inputFile, err)
		}
		if !bytes.Equal(wantData, gotData) {
			t.Errorf("%s: Reversed diff does not match expected.\nWant vs got:\n%s",
				test.inputFile, cmp.Diff(wantData, gotData))
		}
	}
}

func TestReverseRoundTripOnTestdata(t *testing.T) {
	fixtures, err := filepath.Glob(filepath.Join("testdata", "*.diff"))
	if err != nil {
		t.Fatal(err)
	}

	skipped := map[string]bool{
		"empty.diff":            true,
		"empty_new.diff":        true,
		"empty_orig.diff":       true,
		"sample_bad_hunks.diff": true,
	}

	for _, fixture := range fixtures {
		fixture := fixture
		name := filepath.Base(fixture)

		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatal(err)
			}

			if fileDiffs, err := ParseMultiFileDiff(data); err == nil && len(fileDiffs) > 0 {
				if name == "complicated_filenames.diff" {
					if _, err := ReverseMultiFileDiff(fileDiffs); err == nil {
						t.Fatal("reversing fixture with copy diffs succeeded")
					}
					return
				}

				reversed, err := ReverseMultiFileDiff(fileDiffs)
				if err != nil {
					t.Fatalf("first reverse: %s", err)
				}
				roundTrip, err := ReverseMultiFileDiff(reversed)
				if err != nil {
					t.Fatalf("second reverse: %s", err)
				}
				if diff := cmp.Diff(fileDiffs, roundTrip); diff != "" {
					t.Fatalf("double reverse did not restore original file diffs:\n%s", diff)
				}
				return
			}

			if hunks, err := ParseHunks(data); err == nil && len(hunks) > 0 {
				var reversed []*Hunk
				for _, hunk := range hunks {
					inv, err := reverseHunk(hunk)
					if err != nil {
						t.Fatalf("first reverse: %s", err)
					}
					reversed = append(reversed, inv)
				}

				var roundTrip []*Hunk
				for _, hunk := range reversed {
					inv, err := reverseHunk(hunk)
					if err != nil {
						t.Fatalf("second reverse: %s", err)
					}
					roundTrip = append(roundTrip, inv)
				}

				if diff := cmp.Diff(hunks, roundTrip); diff != "" {
					t.Fatalf("double reverse did not restore original hunks:\n%s", diff)
				}
				return
			}

			if skipped[name] {
				return
			}

			t.Fatalf("fixture did not contain parseable file diffs or hunks")
		})
	}
}

func TestReverseFileDiffExtendedHeaders(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "new file",
			input: []string{"diff --git a/f b/f", "new file mode 100644", "index 0000000..587be6b"},
			want:  []string{"diff --git b/f a/f", "deleted file mode 100644", "index 587be6b..0000000"},
		},
		{
			name:  "deleted file",
			input: []string{"diff --git a/f b/f", "deleted file mode 100644", "index 587be6b..0000000"},
			want:  []string{"diff --git b/f a/f", "new file mode 100644", "index 0000000..587be6b"},
		},
		{
			name:  "rename",
			input: []string{"diff --git a/old b/new", "similarity index 70%", "rename from old", "rename to new", "index 94954ab..8b14c4f 100644"},
			want:  []string{"diff --git b/new a/old", "similarity index 70%", "rename from new", "rename to old", "index 8b14c4f..94954ab 100644"},
		},
		{
			name:  "mode change",
			input: []string{"diff --git a/f b/f", "old mode 100644", "new mode 100755"},
			want:  []string{"diff --git b/f a/f", "old mode 100755", "new mode 100644"},
		},
		{
			name:  "CRLF index",
			input: []string{"diff --git a/f b/f\r", "index 94954ab..8b14c4f 100644\r"},
			want:  []string{"diff --git b/f a/f\r", "index 8b14c4f..94954ab 100644\r"},
		},
		{
			name:  "CRLF mode and rename",
			input: []string{"diff --git a/old b/new\r", "old mode 100644\r", "new mode 100755\r", "rename from old\r", "rename to new\r"},
			want:  []string{"diff --git b/new a/old\r", "old mode 100755\r", "new mode 100644\r", "rename from new\r", "rename to old\r"},
		},
		{
			name:  "unknown and malformed headers",
			input: []string{"diff --git a/f b/f", "x-header value", "old mode 100644", "copy from", "index missing-separator"},
			want:  []string{"diff --git b/f a/f", "x-header value", "old mode 100644", "copy from", "index missing-separator"},
		},
		{
			name:  "no extended headers",
			input: nil,
			want:  nil,
		},
		{
			// Only git emits extended headers, so leave anything else alone.
			name:  "non-git header block",
			input: []string{"diff --ruN a/f b/f", "old mode 0777", "new mode 0755"},
			want:  []string{"diff --ruN a/f b/f", "old mode 0777", "new mode 0755"},
		},
	}
	for _, test := range tests {
		orig := append([]string(nil), test.input...)
		fd := &FileDiff{OrigName: "a/f", NewName: "b/f", Extended: test.input}
		reversed, err := ReverseFileDiff(fd)
		if err != nil {
			t.Errorf("%s: ReverseFileDiff: %s", test.name, err)
			continue
		}
		if d := cmp.Diff(test.want, reversed.Extended); d != "" {
			t.Errorf("%s: reversed extended headers differ (-want +got):\n%s", test.name, d)
		}
		if d := cmp.Diff(orig, fd.Extended); d != "" {
			t.Errorf("%s: ReverseFileDiff mutated its input (-want +got):\n%s", test.name, d)
		}
	}
}

func TestReverseFileDiffGitHeader(t *testing.T) {
	// These expected lines were captured from git diff -R. Git is deliberately
	// not invoked by the test.
	tests := []struct {
		name             string
		input, orig, new string
		want             string
	}{
		{
			name:  "unquoted paths with spaces disambiguated by parsed names",
			input: "diff --git a/old name.txt b/new name.txt",
			orig:  "a/old name.txt",
			new:   "b/new name.txt",
			want:  "diff --git b/new name.txt a/old name.txt",
		},
		{
			name:  "quoted newline",
			input: `diff --git "a/line\nbreak.txt" "b/other\nline.txt"`,
			want:  `diff --git "b/other\nline.txt" "a/line\nbreak.txt"`,
		},
		{
			name:  "quoted escapes and non-ASCII",
			input: `diff --git "a/quote\"slash\\\303\251.txt" "b/renamed \"path\"\\\303\270.txt"`,
			want:  `diff --git "b/renamed \"path\"\\\303\270.txt" "a/quote\"slash\\\303\251.txt"`,
		},
		{
			name:  "CRLF",
			input: "diff --git a/plain-old.txt b/plain-new.txt\r",
			want:  "diff --git b/plain-new.txt a/plain-old.txt\r",
		},
		{
			name:  "ambiguous paths without corroborating names",
			input: "diff --git a/old name.txt b/new name.txt",
			want:  "diff --git a/old name.txt b/new name.txt",
		},
		{
			name:  "malformed quoting",
			input: `diff --git "a/old b/new`,
			orig:  "a/old",
			new:   "b/new",
			want:  `diff --git "a/old b/new`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fd := &FileDiff{OrigName: test.orig, NewName: test.new, Extended: []string{test.input}}
			reversed, err := ReverseFileDiff(fd)
			if err != nil {
				t.Fatal(err)
			}
			if got := reversed.Extended[0]; got != test.want {
				t.Errorf("diff --git header: got %q, want %q", got, test.want)
			}
		})
	}
}

func TestReverseFileDiffRejectsCopy(t *testing.T) {
	tests := [][]string{
		{"diff --git a/old b/new", "similarity index 100%", "copy from old", "copy to new"},
		{"diff --git a/old b/new", "copy from old"},
	}
	for _, extended := range tests {
		fd := &FileDiff{Extended: extended}
		if _, err := ReverseFileDiff(fd); err != ErrCannotReverseCopy {
			t.Errorf("ReverseFileDiff error = %v, want ErrCannotReverseCopy for headers %q", err, extended)
		}
	}
}

func TestReverseFileDiffEmptyNewFile(t *testing.T) {
	input := []byte("diff --git a/empty.txt b/empty.txt\nnew file mode 100644\nindex 0000000..e69de29\n")
	fd, err := ParseFileDiff(input)
	if err != nil {
		t.Fatal(err)
	}
	reversed, err := ReverseFileDiff(fd)
	if err != nil {
		t.Fatal(err)
	}
	printed, err := PrintFileDiff(reversed)
	if err != nil {
		t.Fatal(err)
	}
	roundTrip, err := ParseFileDiff(printed)
	if err != nil {
		t.Fatal(err)
	}
	// The forward diff creates the file, so it parses as OrigName=/dev/null.
	// Reversing it must delete the file, i.e. NewName=/dev/null.
	if fd.OrigName != "/dev/null" {
		t.Fatalf("forward diff: got OrigName=%q, want /dev/null", fd.OrigName)
	}
	if roundTrip.NewName != "/dev/null" || roundTrip.OrigName == "/dev/null" {
		t.Errorf("reversed empty new-file diff: got OrigName=%q NewName=%q, want NewName=/dev/null\nprinted:\n%s",
			roundTrip.OrigName, roundTrip.NewName, printed)
	}
}
