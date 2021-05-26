package diff

import (
	"testing"
)

func TestReadQuotedFilename_Success(t *testing.T) {
	tests := []struct {
		input, value, remainder string
	}{
		{input: `""`, value: "", remainder: ""},
		{input: `"aaa"`, value: "aaa", remainder: ""},
		{input: `"aaa" bbb`, value: "aaa", remainder: " bbb"},
		{input: `"aaa" "bbb" ccc`, value: "aaa", remainder: ` "bbb" ccc`},
		{input: `"\""`, value: "\"", remainder: ""},
		{input: `"uh \"oh\""`, value: "uh \"oh\"", remainder: ""},
		{input: `"uh \\"oh\\""`, value: "uh \\", remainder: `oh\\""`},
		{input: `"uh \\\"oh\\\""`, value: "uh \\\"oh\\\"", remainder: ""},
	}
	for _, tc := range tests {
		value, remainder, err := readQuotedFilename(tc.input)
		if err != nil {
			t.Errorf("readQuotedFilename(`%s`): expected success, got '%s'", tc.input, err)
		} else if value != tc.value || remainder != tc.remainder {
			t.Errorf("readQuotedFilename(`%s`): expected `%s` and `%s`, got `%s` and `%s`", tc.input, tc.value, tc.remainder, value, remainder)
		}
	}
}

func TestReadQuotedFilename_Error(t *testing.T) {
	tests := []string{
		// Doesn't start with a quote
		``,
		`foo`,
		` "foo"`,
		// Missing end quote
		`"`,
		`"\"`,
		// "\x" is not a valid Go string literal escape
		`"\xxx"`,
	}
	for _, input := range tests {
		_, _, err := readQuotedFilename(input)
		if err == nil {
			t.Errorf("readQuotedFilename(`%s`): expected error", input)
		}
	}
}

func TestParseDiffGitArgs_Success(t *testing.T) {
	tests := []struct {
		input, first, second string
	}{
		{input: `aaa bbb`, first: "aaa", second: "bbb"},
		{input: `"aaa" bbb`, first: "aaa", second: "bbb"},
		{input: `aaa "bbb"`, first: "aaa", second: "bbb"},
		{input: `"aaa" "bbb"`, first: "aaa", second: "bbb"},
		{input: `1/a 2/z`, first: "1/a", second: "2/z"},
		{input: `1/hello world 2/hello world`, first: "1/hello world", second: "2/hello world"},
		{input: `"new\nline" and spaces`, first: "new\nline", second: "and spaces"},
		{input: `a/existing file with spaces "b/new, complicated\nfilen\303\270me"`, first: "a/existing file with spaces", second: "b/new, complicated\nfilen\303\270me"},
	}
	for _, tc := range tests {
		first, second, success := parseDiffGitArgs(tc.input)
		if !success {
			t.Errorf("`diff --git %s`: expected success", tc.input)
		} else if first != tc.first || second != tc.second {
			t.Errorf("`diff --git %s`: expected `%s` and `%s`, got `%s` and `%s`", tc.input, tc.first, tc.second, first, second)
		}
	}
}

func TestParseDiffGitArgs_Unsuccessful(t *testing.T) {
	tests := []string{
		``,
		`hello_world.txt`,
		`word `,
		` word`,
		`"a/bad_quoting b/bad_quoting`,
		`a/bad_quoting "b/bad_quoting`,
		`a/bad_quoting b/bad_quoting"`,
		`"a/bad_quoting b/bad_quoting"`,
		`"a/bad""b/bad"`,
		`"a/bad" "b/bad" "c/bad"`,
		`a/bad "b/bad" "c/bad"`,
	}
	for _, input := range tests {
		first, second, success := parseDiffGitArgs(input)
		if success {
			t.Errorf("`diff --git %s`: expected unsuccessful; got `%s` and `%s`", input, first, second)
		}
	}
}
