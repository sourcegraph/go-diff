package diff

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func unix(sec int64) *time.Time {
	t := time.Unix(sec, 0)
	return &t
}

func init() {
	// Diffs include times that by default are generated in the local
	// timezone. To ensure that tests behave the same in all timezones
	// (compared to the hard-coded expected output), force the test
	// timezone to UTC.
	//
	// This is safe to do in tests but should not (and need not) be
	// done for the main code.
	time.Local = time.UTC
}

func TestParseHunkNoChunksize(t *testing.T) {
	filename := "sample_no_chunksize.diff"
	diffData, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatal(err)
	}
	diff, err := ParseHunks(diffData)
	if err != nil {
		t.Errorf("%s: got ParseHunks err %v,  want %v", filename, err, nil)
	}
	if len(diff) != 1 {
		t.Fatalf("%s: Got %d hunks, want only one", filename, len(diff))
	}

	correct := &Hunk{
		NewLines:      1,
		NewStartLine:  1,
		OrigLines:     0,
		OrigStartLine: 0,
		StartPosition: 1,
	}
	h := diff[0]
	h.Body = nil // We're not testing the body.
	if !cmp.Equal(h, correct) {
		t.Errorf("%s: got - want:\n%s", filename, cmp.Diff(correct, h))
	}
}

func TestParseHunksAndPrintHunks(t *testing.T) {
	tests := []struct {
		filename     string
		wantParseErr error
	}{
		{filename: "sample_hunk.diff"},
		{filename: "sample_hunks.diff"},
		{filename: "sample_bad_hunks.diff"},
		{filename: "sample_hunks_no_newline.diff"},
		{filename: "no_newline_both.diff"},
		{filename: "no_newline_both2.diff"},
		{filename: "no_newline_orig.diff"},
		{filename: "no_newline_new.diff"},
		{filename: "empty_orig.diff"},
		{filename: "empty_new.diff"},
		{filename: "oneline_hunk.diff"},
		{filename: "empty.diff"},
		{filename: "sample_hunk_lines_start_with_minuses.diff"},
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
			t.Errorf("%s: printed diff hunks != original diff hunks\n\n# PrintHunks output - Original:\n%s", test.filename, cmp.Diff(diffData, printed))
		}
	}
}

func TestParseFileDiffHeaders(t *testing.T) {
	tests := []struct {
		filename string
		wantDiff *FileDiff
	}{
		{
			filename: "sample_file.diff",
			wantDiff: &FileDiff{
				OrigName: "oldname",
				OrigTime: unix(1255273940), // 2009-10-11 15:12:20
				NewName:  "newname",
				NewTime:  unix(1255273950), // 2009-10-11 15:12:30
			},
		},
		{
			filename: "sample_file_no_fractional_seconds.diff",
			wantDiff: &FileDiff{
				OrigName: "goyaml.go",
				OrigTime: unix(1322164040), // 2011-11-24 19:47:20
				NewName:  "goyaml.go",
				NewTime:  unix(1322486679), // 2011-11-28 13:24:39
			},
		},
		{
			filename: "sample_file_extended.diff",
			wantDiff: &FileDiff{
				OrigName: "oldname",
				OrigTime: unix(1255273940), // 2009-10-11 15:12:20
				NewName:  "newname",
				NewTime:  unix(1255273950), // 2009-10-11 15:12:30
				Extended: []string{
					"diff --git a/vcs/git_cmd.go b/vcs/git_cmd.go",
					"index aa4de15..7c048ab 100644",
				},
			},
		},
		{
			filename: "sample_file_extended_empty_new.diff",
			wantDiff: &FileDiff{
				OrigName: "/dev/null",
				OrigTime: nil,
				NewName:  "b/vendor/go/build/testdata/empty/dummy",
				NewTime:  nil,
				Extended: []string{
					"diff --git a/vendor/go/build/testdata/empty/dummy b/vendor/go/build/testdata/empty/dummy",
					"new file mode 100644",
					"index 0000000..e69de29",
				},
			},
		},
		{
			filename: "sample_file_extended_empty_new_binary.diff",
			wantDiff: &FileDiff{
				OrigName: "/dev/null",
				OrigTime: nil,
				NewName:  "b/diff/binary-image.png",
				NewTime:  nil,
				Extended: []string{
					"diff --git a/diff/binary-image.png b/diff/binary-image.png",
					"new file mode 100644",
					"index 0000000..b51756e",
					"Binary files /dev/null and b/diff/binary-image.png differ",
				},
			},
		},
		{
			filename: "sample_file_extended_empty_deleted.diff",
			wantDiff: &FileDiff{
				OrigName: "a/vendor/go/build/testdata/empty/dummy",
				OrigTime: nil,
				NewName:  "/dev/null",
				NewTime:  nil,
				Extended: []string{
					"diff --git a/vendor/go/build/testdata/empty/dummy b/vendor/go/build/testdata/empty/dummy",
					"deleted file mode 100644",
					"index e69de29..0000000",
				},
			},
		},
		{
			filename: "sample_file_extended_empty_deleted_binary.diff",
			wantDiff: &FileDiff{
				OrigName: "a/187/player/random/gopher-0.png",
				OrigTime: nil,
				NewName:  "/dev/null",
				NewTime:  nil,
				Extended: []string{
					"diff --git a/187/player/random/gopher-0.png b/187/player/random/gopher-0.png",
					"deleted file mode 100644",
					"index aebdfc7..0000000",
					"Binary files a/187/player/random/gopher-0.png and /dev/null differ",
				},
			},
		},
		{
			filename: "sample_file_extended_empty_rename.diff",
			wantDiff: &FileDiff{
				OrigName: "a/docs/integrations/Email_Notifications.md",
				OrigTime: nil,
				NewName:  "b/docs/integrations/email-notifications.md",
				NewTime:  nil,
				Extended: []string{
					"diff --git a/docs/integrations/Email_Notifications.md b/docs/integrations/email-notifications.md",
					"similarity index 100%",
					"rename from docs/integrations/Email_Notifications.md",
					"rename to docs/integrations/email-notifications.md",
				},
			},
		},
		{
			filename: "quoted_filename.diff",
			wantDiff: &FileDiff{
				OrigName: "a/商品详情.txt",
				OrigTime: nil,
				NewName:  "b/商品详情.txt",
				NewTime:  nil,
				Extended: []string{
					"diff --git \"a/\\345\\225\\206\\345\\223\\201\\350\\257\\246\\346\\203\\205.txt\" \"b/\\345\\225\\206\\345\\223\\201\\350\\257\\246\\346\\203\\205.txt\"",
					"index e69de29..c67479b 100644",
				},
			},
		},
		{
			filename: "sample_file_extended_binary_rename.diff",
			wantDiff: &FileDiff{
				OrigName: "a/data/Font.png",
				OrigTime: nil,
				NewName:  "b/data/Other.png",
				NewTime:  nil,
				Extended: []string{
					"diff --git a/data/Font.png b/data/Other.png",
					"similarity index 51%",
					"rename from data/Font.png",
					"rename to data/Other.png",
					"index 17a971d..599f8dd 100644",
					"Binary files a/data/Font.png and b/data/Other.png differ",
				},
			},
		},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diff, err := ParseFileDiff(diffData)
		if err != nil {
			t.Fatalf("%s: got ParseFileDiff error %v", test.filename, err)
		}

		diff.Hunks = nil
		if got, want := diff, test.wantDiff; !cmp.Equal(got, want) {
			t.Errorf("%s:\n\ngot - want:\n%s", test.filename, cmp.Diff(want, got))
		}
	}
}

func TestParseMultiFileDiffHeaders(t *testing.T) {
	tests := []struct {
		filename  string
		wantDiffs []*FileDiff
	}{
		{
			filename: "sample_multi_file_new.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "/dev/null",
					OrigTime: nil,
					NewName:  "b/_vendor/go/build/syslist_test.go",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/_vendor/go/build/syslist_test.go b/_vendor/go/build/syslist_test.go",
						"new file mode 100644",
						"index 0000000..3be2928",
					},
				},
				{
					OrigName: "/dev/null",
					OrigTime: nil,
					NewName:  "b/_vendor/go/build/testdata/empty/dummy",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/_vendor/go/build/testdata/empty/dummy b/_vendor/go/build/testdata/empty/dummy",
						"new file mode 100644",
						"index 0000000..e69de29",
					},
				},
				{
					OrigName: "/dev/null",
					OrigTime: nil,
					NewName:  "b/_vendor/go/build/testdata/multi/file.go",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/_vendor/go/build/testdata/multi/file.go b/_vendor/go/build/testdata/multi/file.go",
						"new file mode 100644",
						"index 0000000..ee946eb",
					},
				},
			},
		},
		{
			filename: "sample_multi_file_deleted.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "a/vendor/go/build/syslist_test.go",
					OrigTime: nil,
					NewName:  "/dev/null",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/vendor/go/build/syslist_test.go b/vendor/go/build/syslist_test.go",
						"deleted file mode 100644",
						"index 3be2928..0000000",
					},
				},
				{
					OrigName: "a/vendor/go/build/testdata/empty/dummy",
					OrigTime: nil,
					NewName:  "/dev/null",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/vendor/go/build/testdata/empty/dummy b/vendor/go/build/testdata/empty/dummy",
						"deleted file mode 100644",
						"index e69de29..0000000",
					},
				},
				{
					OrigName: "a/vendor/go/build/testdata/multi/file.go",
					OrigTime: nil,
					NewName:  "/dev/null",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/vendor/go/build/testdata/multi/file.go b/vendor/go/build/testdata/multi/file.go",
						"deleted file mode 100644",
						"index ee946eb..0000000",
					},
				},
			},
		},
		{
			filename: "sample_multi_file_rename.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "a/README.md",
					OrigTime: nil,
					NewName:  "b/README.md",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/README.md b/README.md",
						"index 5f3d591..96a24fa 100644",
					},
				},
				{
					OrigName: "a/docs/integrations/Email_Notifications.md",
					OrigTime: nil,
					NewName:  "b/docs/integrations/email-notifications.md",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/docs/integrations/Email_Notifications.md b/docs/integrations/email-notifications.md",
						"similarity index 100%",
						"rename from docs/integrations/Email_Notifications.md",
						"rename to docs/integrations/email-notifications.md",
					},
				},
				{
					OrigName: "a/release_notes.md",
					OrigTime: nil,
					NewName:  "b/release_notes.md",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/release_notes.md b/release_notes.md",
						"index f2ff13f..f060cb5 100644",
					},
				},
			},
		},
		{
			filename: "sample_multi_file_binary.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "a/README.md",
					OrigTime: nil,
					NewName:  "b/README.md",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/README.md b/README.md",
						"index 7b73e04..36cde13 100644",
					},
				},
				{
					OrigName: "a/data/Font.png",
					OrigTime: nil,
					NewName:  "b/data/Font.png",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/data/Font.png b/data/Font.png",
						"index 17a971d..599f8dd 100644",
						"Binary files a/data/Font.png and b/data/Font.png differ",
					},
				},
				{
					OrigName: "a/main.go",
					OrigTime: nil,
					NewName:  "b/main.go",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/main.go b/main.go",
						"index 1aced1e..98a982e 100644",
					},
				},
			},
		},
		{
			filename: "sample_binary_inline.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "a/logo-old.png",
					OrigTime: nil,
					NewName:  "/dev/null",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/logo-old.png b/logo-old.png",
						"deleted file mode 100644",
						"index d29d0e9757e0d9b854a8ed58f170bcb454cc1ae3..0000000000000000000000000000000000000000",
						"GIT binary patch",
						"literal 0",
						"HcmV?d00001",
						"",
						"literal 0",
						"HcmV?d00001",
						"",
					},
				},
				{
					OrigName: "a/logo-old.png",
					OrigTime: nil,
					NewName:  "b/logo-old.png",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/logo-old.png b/logo-old.png",
						"index ff82e793467f2050d731d75b4968d2e6b9c5d33b..d29d0e9757e0d9b854a8ed58f170bcb454cc1ae3 100644",
						"GIT binary patch",
						"literal 0",
						"HcmV?d00001",
						"",
						"literal 0",
						"HcmV?d00001",
						"",
					},
				},
				{
					OrigName: "a/logo.png",
					OrigTime: nil,
					NewName:  "b/logo-old.png",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/logo.png b/logo-old.png",
						"similarity index 100%",
						"rename from logo.png",
						"rename to logo-old.png",
					},
				},
				{
					OrigName: "/dev/null",
					OrigTime: nil,
					NewName:  "b/logo.png",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/logo.png b/logo.png",
						"new file mode 100644",
						"index 0000000000000000000000000000000000000000..ff82e793467f2050d731d75b4968d2e6b9c5d33b",
						"GIT binary patch",
						"literal 0",
						"HcmV?d00001",
						"",
						"literal 0",
						"HcmV?d00001",
						"",
					},
				},
			},
		},
		{
			filename: "sample_multi_file_new_win.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "/dev/null",
					OrigTime: nil,
					NewName:  "b/_vendor/go/build/syslist_test.go",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/_vendor/go/build/syslist_test.go b/_vendor/go/build/syslist_test.go",
						"new file mode 100644",
						"index 0000000..3be2928",
					},
				},
				{
					OrigName: "/dev/null",
					OrigTime: nil,
					NewName:  "b/_vendor/go/build/testdata/empty/dummy",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/_vendor/go/build/testdata/empty/dummy b/_vendor/go/build/testdata/empty/dummy",
						"new file mode 100644",
						"index 0000000..e69de29",
					},
				},
				{
					OrigName: "/dev/null",
					OrigTime: nil,
					NewName:  "b/_vendor/go/build/testdata/multi/file.go",
					NewTime:  nil,
					Extended: []string{
						"diff --git a/_vendor/go/build/testdata/multi/file.go b/_vendor/go/build/testdata/multi/file.go",
						"new file mode 100644",
						"index 0000000..ee946eb",
					},
				},
			},
		},
		{
			filename: "sample_contains_added_deleted_files.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "source_a/file_1.txt",
					OrigTime: nil,
					NewName:  "source_b/file_1.txt",
					NewTime:  nil,
					Extended: []string{
						"diff -u source_a/file_1.txt  source_b/file_1.txt",
					},
				},
				{
					OrigName: "source_a/file_2.txt",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
				{
					OrigName: "source_b/file_3.txt",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
			},
		},
		{
			filename: "sample_contains_only_added_deleted_files.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "source_a/file_1.txt",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
				{
					OrigName: "source_a/file_2.txt",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
				{
					OrigName: "source_b/file_3.txt",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
			},
		},
		{
			filename: "sample_onlyin_line_isnt_a_file_header.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "source_a/file_1.txt",
					OrigTime: nil,
					NewName:  "source_b/file_1.txt",
					NewTime:  nil,
					Extended: []string{
						"diff -u source_a/file_1.txt  source_b/file_1.txt",
					},
				},
				{
					OrigName: "source_a/file_2.txt",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: []string{
						"Only in universe!",
					},
				},
				{
					OrigName: "source_b/file_3.txt some unrelated stuff here.",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
				{
					OrigName: "source_b/file_3.txt",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
			},
		},
		{
			filename: "sample_onlyin_complex_filenames.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "internal/trace/foo bar/bam",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
				{
					OrigName: "internal/trace/foo bar/bam: bar",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
				{
					OrigName: "internal/trace/hello/world: bazz",
					OrigTime: nil,
					NewName:  "",
					NewTime:  nil,
					Extended: nil,
				},
			},
		},
		{
			filename: "sample_multi_file_without_extended.diff",
			wantDiffs: []*FileDiff{
				{
					OrigName: "source_1_a/file_1.txt",
					OrigTime: nil,
					NewName:  "source_1_c/file_1.txt",
					NewTime:  nil,
					Extended: nil,
				},
				{
					OrigName: "source_1_a/file_2.txt",
					OrigTime: nil,
					NewName:  "source_1_c/file_2.txt",
					NewTime:  nil,
					Extended: nil,
				},
			},
		},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diffs, err := ParseMultiFileDiff(diffData)
		if err != nil {
			t.Fatalf("%s: got ParseMultiFileDiff error %v", test.filename, err)
		}

		for i := range diffs {
			diffs[i].Hunks = nil // This test focuses on things other than hunks, so don't compare them.
		}
		if got, want := diffs, test.wantDiffs; !cmp.Equal(got, want) {
			t.Errorf("%s:\n\ngot - want:\n%s", test.filename, cmp.Diff(want, got))
		}
	}
}

func TestParseFileDiffAndPrintFileDiff(t *testing.T) {
	tests := []struct {
		filename     string
		wantParseErr error
	}{
		{filename: "sample_file.diff"},
		{filename: "sample_file_no_timestamp.diff"},
		{filename: "sample_file_extended.diff"},
		{filename: "sample_file_extended_empty_new.diff"},
		{filename: "sample_file_extended_empty_new_binary.diff"},
		{filename: "sample_file_extended_empty_deleted.diff"},
		{filename: "sample_file_extended_empty_deleted_binary.diff"},
		{filename: "sample_file_extended_empty_rename.diff"},
		{filename: "sample_file_extended_empty_binary.diff"},
		{
			filename:     "empty.diff",
			wantParseErr: &ParseError{0, 0, ErrExtendedHeadersEOF},
		},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diff, err := ParseFileDiff(diffData)
		if !reflect.DeepEqual(err, test.wantParseErr) {
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
			t.Errorf("%s: printed file diff != original file diff\n\n# PrintFileDiff output - Original:\n%s", test.filename, cmp.Diff(diffData, printed))
		}
	}
}

func TestParseMultiFileDiffAndPrintMultiFileDiff(t *testing.T) {
	tests := []struct {
		filename      string
		wantParseErr  error
		wantFileDiffs int // How many instances of diff.FileDiff are expected.
	}{
		{filename: "sample_multi_file.diff", wantFileDiffs: 2},
		{filename: "sample_multi_file_single.diff", wantFileDiffs: 1},
		{filename: "sample_multi_file_new.diff", wantFileDiffs: 3},
		{filename: "sample_multi_file_deleted.diff", wantFileDiffs: 3},
		{filename: "sample_multi_file_rename.diff", wantFileDiffs: 3},
		{filename: "sample_multi_file_binary.diff", wantFileDiffs: 3},
		{filename: "long_line_multi.diff", wantFileDiffs: 3},
		{filename: "empty.diff", wantFileDiffs: 0},
		{filename: "empty_multi.diff", wantFileDiffs: 2},
		{filename: "sample_contains_added_deleted_files.diff", wantFileDiffs: 3},
		{filename: "sample_contains_only_added_deleted_files.diff", wantFileDiffs: 3},
		{filename: "sample_onlyin_line_isnt_a_file_header.diff", wantFileDiffs: 4},
		{filename: "sample_onlyin_complex_filenames.diff", wantFileDiffs: 3},
	}
	for _, test := range tests {
		diffData, err := ioutil.ReadFile(filepath.Join("testdata", test.filename))
		if err != nil {
			t.Fatal(err)
		}
		diffs, err := ParseMultiFileDiff(diffData)
		if err != test.wantParseErr {
			t.Errorf("%s: got ParseMultiFileDiff err %v, want %v", test.filename, err, test.wantParseErr)
			continue
		}
		if test.wantParseErr != nil {
			continue
		}

		if got, want := len(diffs), test.wantFileDiffs; got != want {
			t.Errorf("%s: got %v instances of diff.FileDiff, expected %v", test.filename, got, want)
		}

		printed, err := PrintMultiFileDiff(diffs)
		if err != nil {
			t.Errorf("%s: PrintMultiFileDiff: %s", test.filename, err)
		}
		if !bytes.Equal(printed, diffData) {
			t.Errorf("%s: printed multi-file diff != original multi-file diff\n\n# PrintMultiFileDiff output - Original:\n%s", test.filename, cmp.Diff(diffData, printed))
		}
	}
}

func TestNoNewlineAtEnd(t *testing.T) {
	diffs := map[string]struct {
		diff              string
		trailingNewlineOK bool
	}{
		"orig": {
			diff: `@@ -1,1 +1,1 @@
-a
\ No newline at end of file
+b
`,
			trailingNewlineOK: true,
		},
		"new": {
			diff: `@@ -1,1 +1,1 @@
-a
+b
\ No newline at end of file
`,
		},
		"both": {
			diff: `@@ -1,1 +1,1 @@
-a
\ No newline at end of file
+b
\ No newline at end of file
`,
		},
	}

	for label, test := range diffs {
		hunks, err := ParseHunks([]byte(test.diff))
		if err != nil {
			t.Errorf("%s: ParseHunks: %s", label, err)
			continue
		}

		for _, hunk := range hunks {
			if body := string(hunk.Body); strings.Contains(body, "No newline") {
				t.Errorf("%s: after parse, hunk body contains 'No newline...' string\n\nbody is:\n%s", label, body)
			}
			if !test.trailingNewlineOK {
				if bytes.HasSuffix(hunk.Body, []byte{'\n'}) {
					t.Errorf("%s: after parse, hunk body ends with newline\n\nbody is:\n%s", label, hunk.Body)
				}
			}
			if dontWant := []byte("-a+b"); bytes.Contains(hunk.Body, dontWant) {
				t.Errorf("%s: hunk body contains %q\n\nbody is:\n%s", label, dontWant, hunk.Body)
			}

			printed, err := PrintHunks(hunks)
			if err != nil {
				t.Errorf("%s: PrintHunks: %s", label, err)
				continue
			}
			if printed := string(printed); printed != test.diff {
				t.Errorf("%s: printed diff hunks != original diff hunks\n\n# PrintHunks output - Original:\n%s", label, cmp.Diff(test.diff, printed))
			}
		}
	}
}

func TestFileDiff_Stat(t *testing.T) {
	tests := map[string]struct {
		hunks []*Hunk
		want  Stat
	}{
		"no change": {
			hunks: []*Hunk{
				{Body: []byte(`@@ -0,0 +0,0
 a
 b
`)},
			},
			want: Stat{},
		},
		"added/deleted": {
			hunks: []*Hunk{
				{Body: []byte(`@@ -0,0 +0,0
+a
 b
-c
 d
`)},
			},
			want: Stat{Added: 1, Deleted: 1},
		},
		"changed": {
			hunks: []*Hunk{
				{Body: []byte(`@@ -0,0 +0,0
+a
+b
-c
-d
 e
`)},
			},
			want: Stat{Added: 1, Changed: 1, Deleted: 1},
		},
		"many changes": {
			hunks: []*Hunk{
				{Body: []byte(`@@ -0,0 +0,0
+a
-b
+c
-d
 e
`)},
			},
			want: Stat{Added: 0, Changed: 2, Deleted: 0},
		},
	}
	for label, test := range tests {
		fdiff := &FileDiff{Hunks: test.hunks}
		stat := fdiff.Stat()
		if !cmp.Equal(stat, test.want) {
			t.Errorf("%s: got - want diff stat\n%s", label, cmp.Diff(test.want, stat))
			continue
		}
	}
}
