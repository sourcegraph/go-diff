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
