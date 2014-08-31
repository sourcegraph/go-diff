package diff

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseAndPrint(t *testing.T) {
	tests := []struct {
		filename     string
		wantParseErr error
	}{
		{
			filename: "sample_single_hunk.diff",
		},
		{
			filename: "sample_single_file.diff",
		},
		{
			filename:     "sample_bad.diff",
			wantParseErr: nil,
		},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diff, err := Parse(diffData)
		if err != test.wantParseErr {
			t.Errorf("%s: got Parse err %v, want %v", test.filename, err, test.wantParseErr)
			continue
		}
		if test.wantParseErr != nil {
			continue
		}

		printedDiff, err := Print(diff)
		if err != nil {
			t.Errorf("%s: Print: %s", test.filename, err)
		}
		if !bytes.Equal(printedDiff, diffData) {
			t.Errorf("%s: printed diff != original diff\n\n# Print output:\n%s\n\n# Original:\n%s", test.filename, printedDiff, diffData)
		}
	}
}
