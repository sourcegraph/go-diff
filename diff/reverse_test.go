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
