package diff

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseHunksAndPrintHunks(t *testing.T) {
	tests := []struct {
		filename     string
		wantParseErr error
	}{
		{
			filename: "sample_hunk.diff",
		},
		{
			filename: "sample_hunks.diff",
		},
		{
			filename:     "sample_bad_hunks.diff",
			wantParseErr: nil,
		},
		{filename: "empty.diff"},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diff, err := ParseHunks(diffData)
		if err != test.wantParseErr {
			t.Errorf("%s: got ParseHunks err %v, want %v", test.filename, err, test.wantParseErr)
			continue
		}
		if test.wantParseErr != nil {
			continue
		}

		printed, err := PrintHunks(diff)
		if err != nil {
			t.Errorf("%s: PrintHunks: %s", test.filename, err)
		}
		if !bytes.Equal(printed, diffData) {
			t.Errorf("%s: printed diff hunks != original diff hunks\n\n# PrintHunks output:\n%s\n\n# Original:\n%s", test.filename, printed, diffData)
		}
	}
}

func TestParseFileDiffAndPrintFileDiff(t *testing.T) {
	tests := []struct {
		filename     string
		wantParseErr error
	}{
		{
			filename: "sample_file.diff",
		},
		{
			filename: "sample_file_no_timestamp.diff",
		},
		{
			filename: "sample_file_extended.diff",
		},
		{filename: "empty.diff"},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diff, err := ParseFileDiff(diffData)
		if err != test.wantParseErr {
			t.Errorf("%s: got ParseFileDiff err %v, want %v", test.filename, err, test.wantParseErr)
			continue
		}
		if test.wantParseErr != nil {
			continue
		}

		printed, err := PrintFileDiff(diff)
		if err != nil {
			t.Errorf("%s: PrintFileDiff: %s", test.filename, err)
		}
		if !bytes.Equal(printed, diffData) {
			t.Errorf("%s: printed file diff != original file diff\n\n# PrintFileDiff output:\n%s\n\n# Original:\n%s", test.filename, printed, diffData)
		}
	}
}

func TestParseMultiFileDiffAndPrintMultiFileDiff(t *testing.T) {
	tests := []struct {
		filename     string
		wantParseErr error
	}{
		{
			filename: "sample_multi_file.diff",
		},
		{
			filename: "sample_multi_file_single.diff",
		},
		{filename: "empty.diff"},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diff, err := ParseMultiFileDiff(diffData)
		if err != test.wantParseErr {
			t.Errorf("%s: got ParseMultiFileDiff err %v, want %v", test.filename, err, test.wantParseErr)
			continue
		}
		if test.wantParseErr != nil {
			continue
		}

		printed, err := PrintMultiFileDiff(diff)
		if err != nil {
			t.Errorf("%s: PrintMultiFileDiff: %s", test.filename, err)
		}
		if !bytes.Equal(printed, diffData) {
			t.Errorf("%s: printed multi-file diff != original multi-file diff\n\n# PrintMultiFileDiff output:\n%s\n\n# Original:\n%s", test.filename, printed, diffData)
		}
	}
}
