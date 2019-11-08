package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		lineCtx LineContext
		want    LineResult
	}{
		{
			name: "whitespace",
			text: `		`,
			want: LineResult{
				IsWhitespace: true,
			},
		},
		{
			name: "a comment",
			text: `// this is a comment line`,
			want: LineResult{
				IsOnlyComment: true,
			},
		},
		{
			name: "function definition",
			text: `func testFunc(s string) error {`,
			want: LineResult{
				NumBracesDiff: +1,
			},
		},
		{
			name: "function definition with inline comment",
			text: `func testFunc(s string) error { // this is an inline comment`,
			want: LineResult{
				IsInlineComment: true,
				NumBracesDiff:   +1,
			},
		},
		{
			name: "simple err check",
			text: `if err != nil {`,
			want: LineResult{
				StartsErrCheck: true,
				NumBracesDiff:  +1,
			},
		},
		{
			name: "complicated err check",
			text: `if fileErr != nil && fileErr != io.EOF {`,
			want: LineResult{
				StartsErrCheck: true,
				NumBracesDiff:  +1,
			},
		},
		{
			name: "variable declaration 1",
			text: `var isVarOfSomeKind = true && false`,
			want: LineResult{}, // everything should be false
		},
		{
			name: "variable declaration 2",
			text: `isVarOfSomeKind := someFunc(a)`,
			want: LineResult{}, // everything should be false
		},
		{
			name: "starts a block comment",
			text: `/* This line starts a block comment...`,
			want: LineResult{
				IsOnlyComment:      true,
				StartsBlockComment: true,
			},
		},
		{
			name:    "ends a block comment",
			text:    `and this is the last line of a block comment*/`,
			lineCtx: LineContext{InBlockComment: true},
			want: LineResult{
				IsOnlyComment:    true,
				EndsBlockComment: true,
			},
		},
		{
			name: "some code with a start of block comment",
			text: `var a int /* This is a block comment start...`,
			want: LineResult{
				IsInlineComment:    true,
				StartsBlockComment: true,
			},
		},
		{
			name: "starts and stops a block comment",
			text: `/* This line starts a block comment which ends */`,
			want: LineResult{
				IsOnlyComment: true,
			},
		},
		{
			name: "some code with a start and stop of a block comment",
			text: `var a int /* This is a block comment start which ends */`,
			want: LineResult{
				IsInlineComment: true,
			},
		},
		{
			name: "ignore code in a comment",
			text: `// if err = nil {`,
			want: LineResult{
				IsOnlyComment: true,
			},
		},
		{
			name: "ignore stuff in strings",
			text: `var a = "if err = nil { /*...//"`,
			want: LineResult{},
		},
		{
			name: "detect when a multiline string starts",
			text: "var a = ` this is first line of a multiline string...",
			want: LineResult{
				StartBacktickBlock: true,
			},
		},
		{
			name:    "detect when a multiline string is ended",
			text:    "...and hence the string ends`",
			lineCtx: LineContext{InBacktickString: true},
			want: LineResult{
				EndsBacktickBlock: true,
			},
		},
		{
			name:    "detect when a multiline string is ended",
			text:    "...and hence the string ends`",
			lineCtx: LineContext{InBacktickString: true},
			want: LineResult{
				EndsBacktickBlock: true,
			},
		},
		{
			name: "starting a scope",
			text: "if a == b {",
			want: LineResult{
				NumBracesDiff: +1,
			},
		},
		{
			name: "ending a scope",
			text: "}",
			want: LineResult{
				NumBracesDiff: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processLine(tt.text, tt.lineCtx)
			assert.Equal(t, tt.want, got)
		})
	}
}
