package main

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessBufReader(t *testing.T) {

	tests := []struct {
		name    string
		text    string
		want    Results
		wantErr bool
	}{
		{
			name: "sample test file 1",
			text: sampleFileA,
			want: Results{
				NumOfFiles:          1,
				LinesOfCode:         74,
				LinesOfErrCheck:     3,
				LinesOfComments:     4,
				LinesWhitespace:     16,
				TotalLinesProcessed: 96,
				NumInlineComments:   0,
				MaxCurlyBracesDepth: 2,
				MaxCurlyBracesDepthLocation: Location{
					Line: 16,
					File: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buff := bytes.NewBufferString(tt.text)
			reader := bufio.NewReader(buff)
			got, err := processBufReader(reader)
			assert.Equal(t, tt.wantErr, err != nil, "got error: %s", err)
			assert.Equal(t, tt.want, got)

		})
	}
}
