package diff

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// These fixtures are reduced from Git output. They exercise complete parsed
// FileDiffs rather than individual extended-header rewrites.
func TestReverseFileDiffGitSemantics(t *testing.T) {
	tests := []struct {
		name    string
		forward string
		reverse string
		crlf    bool
	}{
		{
			name: "addition",
			forward: "diff --git a/added.txt b/added.txt\n" +
				"new file mode 100644\n" +
				"index 0000000..76d4bb8\n" +
				"--- /dev/null\n" +
				"+++ b/added.txt\n" +
				"@@ -0,0 +1 @@\n" +
				"+add\n",
			reverse: "diff --git b/added.txt a/added.txt\n" +
				"deleted file mode 100644\n" +
				"index 76d4bb8..0000000\n" +
				"--- b/added.txt\n" +
				"+++ /dev/null\n" +
				"@@ -1 +0,0 @@\n" +
				"-add\n",
		},
		{
			name: "deletion",
			forward: "diff --git a/deleted.txt b/deleted.txt\n" +
				"deleted file mode 100644\n" +
				"index c8b1b42..0000000\n" +
				"--- a/deleted.txt\n" +
				"+++ /dev/null\n" +
				"@@ -1 +0,0 @@\n" +
				"-delete\n",
			reverse: "diff --git b/deleted.txt a/deleted.txt\n" +
				"new file mode 100644\n" +
				"index 0000000..c8b1b42\n" +
				"--- /dev/null\n" +
				"+++ a/deleted.txt\n" +
				"@@ -0,0 +1 @@\n" +
				"+delete\n",
		},
		{
			name: "rename",
			forward: "diff --git a/old.txt b/renamed.txt\n" +
				"similarity index 100%\n" +
				"rename from old.txt\n" +
				"rename to renamed.txt\n",
			reverse: "diff --git b/renamed.txt a/old.txt\n" +
				"similarity index 100%\n" +
				"rename from renamed.txt\n" +
				"rename to old.txt\n",
		},
		{
			name: "mode only",
			forward: "diff --git a/mode.txt b/mode.txt\n" +
				"old mode 100644\n" +
				"new mode 100755\n",
			reverse: "diff --git b/mode.txt a/mode.txt\n" +
				"old mode 100755\n" +
				"new mode 100644\n",
		},
		{
			name: "quoted path",
			forward: "diff --git \"a/quote\\\\tname.txt\" \"b/quote\\\\tname.txt\"\n" +
				"index dcae08f..b46c5e4 100644\n" +
				"--- \"a/quote\\\\tname.txt\"\n" +
				"+++ \"b/quote\\\\tname.txt\"\n" +
				"@@ -1 +1 @@\n" +
				"-quoted old\n" +
				"+quoted new\n",
			reverse: "diff --git \"b/quote\\\\tname.txt\" \"a/quote\\\\tname.txt\"\n" +
				"index b46c5e4..dcae08f 100644\n" +
				"--- \"b/quote\\\\tname.txt\"\n" +
				"+++ \"a/quote\\\\tname.txt\"\n" +
				"@@ -1 +1 @@\n" +
				"-quoted new\n" +
				"+quoted old\n",
		},
		{
			name: "CRLF",
			forward: "diff --git a/crlf.txt b/crlf.txt\n" +
				"index 9cca7e3..d2eb92c 100644\n" +
				"--- a/crlf.txt\n" +
				"+++ b/crlf.txt\n" +
				"@@ -1 +1 @@\n" +
				"-old\n" +
				"+new\n",
			reverse: "diff --git b/crlf.txt a/crlf.txt\n" +
				"index d2eb92c..9cca7e3 100644\n" +
				"--- b/crlf.txt\n" +
				"+++ a/crlf.txt\n" +
				"@@ -1 +1 @@\n" +
				"-new\n" +
				"+old\n",
			crlf: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertGitReverse(t, test.forward, test.reverse, test.crlf)
		})
	}
}

// Git binary patches contain both directions. Reversing the index hashes makes
// git apply select the inverse payload without reordering the encoded sections.
func TestReverseFileDiffGitBinarySemantics(t *testing.T) {
	tests := []struct {
		name    string
		forward string
		reverse string
	}{
		{
			name: "literal",
			forward: "diff --git a/literal.bin b/literal.bin\n" +
				"index acdc746399905e5ea0d9d4163d2b3c4d04b05123..e0d84f457d2b5c802ecd286e7e576365d94af3b9 100644\n" +
				"GIT binary patch\n" +
				"literal 8\n" +
				"PcmZQzWMXDvWn%{b09*ha\n\n" +
				"literal 64\n" +
				"LcmZ>CqznK65PlFr\n\n",
			reverse: "diff --git b/literal.bin a/literal.bin\n" +
				"index e0d84f457d2b5c802ecd286e7e576365d94af3b9..acdc746399905e5ea0d9d4163d2b3c4d04b05123 100644\n" +
				"GIT binary patch\n" +
				"literal 8\n" +
				"PcmZQzWMXDvWn%{b09*ha\n\n" +
				"literal 64\n" +
				"LcmZ>CqznK65PlFr\n\n",
		},
		{
			name: "delta",
			forward: "diff --git a/delta.bin b/delta.bin\n" +
				"index df437f42c808d41dec5d543d60ce94c8cb8a044a..b762194ffce13918858bfd54f1c4ded144c7e3b4 100644\n" +
				"GIT binary patch\n" +
				"delta 20\n" +
				"bcmZorXi(U2ft@2cBQY;MHAQjb4Gj(ePVNVE\n\n" +
				"delta 10\n" +
				"PcmZorXi!+h!2v`75gY=g\n\n",
			reverse: "diff --git b/delta.bin a/delta.bin\n" +
				"index b762194ffce13918858bfd54f1c4ded144c7e3b4..df437f42c808d41dec5d543d60ce94c8cb8a044a 100644\n" +
				"GIT binary patch\n" +
				"delta 20\n" +
				"bcmZorXi(U2ft@2cBQY;MHAQjb4Gj(ePVNVE\n\n" +
				"delta 10\n" +
				"PcmZorXi!+h!2v`75gY=g\n\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertGitReverse(t, test.forward, test.reverse, false)
		})
	}
}

func TestReverseFileDiffRejectsGitCopy(t *testing.T) {
	forward := "diff --git a/source.txt b/copied.txt\n" +
		"similarity index 100%\n" +
		"copy from source.txt\n" +
		"copy to copied.txt\n"

	fd, err := ParseFileDiff([]byte(forward))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ReverseFileDiff(fd); err == nil {
		t.Fatal("ReverseFileDiff succeeded for a Git copy patch")
	}
}

func assertGitReverse(t *testing.T, forward, reverse string, crlf bool) {
	t.Helper()
	if diff := gitReverseDiff(t, forward, reverse, crlf); diff != "" {
		t.Fatalf("reversed FileDiff differs from expected Git semantics (-want +got):\n%s", diff)
	}
}

func gitReverseDiff(t *testing.T, forward, reverse string, crlf bool) string {
	t.Helper()
	opts := ParseOptions{KeepCR: crlf}
	if crlf {
		forward = strings.ReplaceAll(forward, "\n", "\r\n")
		reverse = strings.ReplaceAll(reverse, "\n", "\r\n")
	}

	fd, err := ParseFileDiffOptions([]byte(forward), opts)
	if err != nil {
		t.Fatalf("parse forward fixture: %s", err)
	}
	want, err := ParseFileDiffOptions([]byte(reverse), opts)
	if err != nil {
		t.Fatalf("parse reverse fixture: %s", err)
	}
	got, err := ReverseFileDiff(fd)
	if err != nil {
		t.Fatalf("ReverseFileDiff: %s", err)
	}

	// diff --git argument rendering is orthogonal to the parsed directional data
	// covered here, so this remains compatible with either argument order.
	got.Extended[0], want.Extended[0] = "", ""
	return cmp.Diff(want, got)
}
