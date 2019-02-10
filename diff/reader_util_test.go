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
