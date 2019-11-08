package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddResults(t *testing.T) {
	var sampleResultsA = Results{
		NumOfFiles:          2,
		LinesOfCode:         3,
		LinesOfErrCheck:     4,
		LinesOfComments:     5,
		LinesWhitespace:     6,
		TotalLinesProcessed: 7,
		NumInlineComments:   8,
		MaxCurlyBracesDepth: 8,
		MaxCurlyBracesDepthLocation: Location{
			File: "sampleResultsA.go",
			Line: 10,
		},
	}

	var sampleResultsB = Results{
		NumOfFiles:          4,
		LinesOfCode:         6,
		LinesOfErrCheck:     8,
		LinesOfComments:     10,
		LinesWhitespace:     12,
		TotalLinesProcessed: 14,
		NumInlineComments:   16,
		MaxCurlyBracesDepth: 19,
		MaxCurlyBracesDepthLocation: Location{
			File: "sampleResultsB.go",
			Line: 20,
		},
	}

	var sampleResultsAB = Results{
		NumOfFiles:          6,
		LinesOfCode:         9,
		LinesOfErrCheck:     12,
		LinesOfComments:     15,
		LinesWhitespace:     18,
		TotalLinesProcessed: 21,
		NumInlineComments:   24,
		MaxCurlyBracesDepth: 19, // stays same
		MaxCurlyBracesDepthLocation: Location{
			File: "sampleResultsB.go",
			Line: 20,
		},
	}

	tests := []struct {
		name string
		a    Results
		b    Results
		want Results
	}{
		{
			name: "new result",
			a:    Results{},
			b:    sampleResultsA,
			want: sampleResultsA,
		},
		{
			name: "new result",
			a:    sampleResultsA,
			b:    sampleResultsB,
			want: sampleResultsAB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := addResults(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}
