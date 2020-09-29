package diff

import (
	"bufio"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty",
			input: "",
			want:  []string{},
		},
		{
			name:  "single_line",
			input: "@@ -0,0 +1,62 @@",
			want:  []string{"@@ -0,0 +1,62 @@"},
		},
		{
			name:  "single_lf_terminated_line",
			input: "@@ -0,0 +1,62 @@\n",
			want:  []string{"@@ -0,0 +1,62 @@"},
		},
		{
			name:  "single_crlf_terminated_line",
			input: "@@ -0,0 +1,62 @@\r\n",
			want:  []string{"@@ -0,0 +1,62 @@"},
		},
		{
			name: "multi_line",
			input: `diff --git a/test.go b/test.go
new file mode 100644
index 0000000..3be2928`,
			want: []string{
				"diff --git a/test.go b/test.go",
				"new file mode 100644",
				"index 0000000..3be2928",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := bufio.NewReader(strings.NewReader(test.input))
			out := []string{}
			for {
				l, err := readLine(in)
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				out = append(out, string(l))
			}
			if !reflect.DeepEqual(test.want, out) {
				t.Errorf("read lines not equal: want %v, got %v", test.want, out)
			}
		})
	}
}

func TestLineReader_ReadLine(t *testing.T) {
	input := `diff --git a/test.go b/test.go
new file mode 100644
index 0000000..3be2928


`

	in := newLineReader(strings.NewReader(input))
	out := []string{}
	for i := 0; i < 4; i++ {
		l, err := in.readLine()
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, string(l))
	}

	wantOut := strings.Split(input, "\n")[0:4]
	if !reflect.DeepEqual(wantOut, out) {
		t.Errorf("read lines not equal: want %v, got %v", wantOut, out)
	}

	_, err := in.readLine()
	if err != nil {
		t.Fatal(err)
	}
	if in.cachedNextLineErr != io.EOF {
		t.Fatalf("lineReader has wrong cachedNextLineErr: %s", in.cachedNextLineErr)
	}
	_, err = in.readLine()
	if err != io.EOF {
		t.Fatalf("readLine did not return io.EOF: %s", err)
	}
}

func TestLineReader_NextLine(t *testing.T) {
	input := `aaa rest of line
bbbrest of line
ccc rest of line`

	in := newLineReader(strings.NewReader(input))

	type assertion struct {
		prefix string
		want   bool
	}

	testsPerReadLine := []struct {
		nextLine        []assertion
		nextNextLine    []assertion
		wantReadLineErr error
	}{
		{
			nextLine: []assertion{
				{prefix: "a", want: true},
				{prefix: "aa", want: true},
				{prefix: "aaa", want: true},
				{prefix: "bbb", want: false},
				{prefix: "ccc", want: false},
			},
			nextNextLine: []assertion{
				{prefix: "aaa", want: false},
				{prefix: "bbb", want: true},
				{prefix: "ccc", want: false},
			},
		},
		{
			nextLine: []assertion{
				{prefix: "aaa", want: false},
				{prefix: "bbb", want: true},
				{prefix: "ccc", want: false},
			},
			nextNextLine: []assertion{
				{prefix: "aaa", want: false},
				{prefix: "bbb", want: false},
				{prefix: "ccc", want: true},
			},
		},
		{
			nextLine: []assertion{
				{prefix: "aaa", want: false},
				{prefix: "bbb", want: false},
				{prefix: "ccc", want: true},
				{prefix: "ddd", want: false},
			},
			nextNextLine: []assertion{
				{prefix: "aaa", want: false},
				{prefix: "bbb", want: false},
				{prefix: "ccc", want: false},
				{prefix: "ddd", want: false},
			},
		},
		{
			nextLine: []assertion{
				{prefix: "aaa", want: false},
				{prefix: "bbb", want: false},
				{prefix: "ccc", want: false},
				{prefix: "ddd", want: false},
			},
			nextNextLine: []assertion{
				{prefix: "aaa", want: false},
				{prefix: "bbb", want: false},
				{prefix: "ccc", want: false},
				{prefix: "ddd", want: false},
			},
			wantReadLineErr: io.EOF,
		},
	}

	for _, tc := range testsPerReadLine {
		for _, assert := range tc.nextLine {
			got, err := in.nextLineStartsWith(assert.prefix)
			if err != nil {
				t.Fatalf("nextLineStartsWith returned unexpected error: %s", err)
			}

			if got != assert.want {
				t.Fatalf("unexpected result for prefix %q. got=%t, want=%t", assert.prefix, got, assert.want)
			}
		}

		for _, assert := range tc.nextNextLine {
			got, err := in.nextNextLineStartsWith(assert.prefix)
			if err != nil {
				t.Fatalf("nextLineStartsWith returned unexpected error: %s", err)
			}

			if got != assert.want {
				t.Fatalf("unexpected result for prefix %q. got=%t, want=%t", assert.prefix, got, assert.want)
			}
		}

		_, err := in.readLine()
		if err != tc.wantReadLineErr {
			t.Fatalf("readLine returned unexpected error. got=%s, want=%s", err, tc.wantReadLineErr)
		}

	}
}
