package main

type Results struct {
	NumOfFiles                  int
	LinesOfCode                 int
	LinesOfErrCheck             int
	LinesOfComments             int
	LinesWhitespace             int
	TotalLinesProcessed         int
	NumInlineComments           int
	MaxCurlyBracesDepth         int
	MaxCurlyBracesDepthLocation Location
}

type Location struct {
	File string
	Line int
}

func addResults(a, b Results) Results {
	var r Results
	r.NumOfFiles = a.NumOfFiles + b.NumOfFiles
	r.LinesOfCode = a.LinesOfCode + b.LinesOfCode
	r.LinesOfComments = a.LinesOfComments + b.LinesOfComments
	r.LinesOfErrCheck = a.LinesOfErrCheck + b.LinesOfErrCheck
	r.LinesWhitespace = a.LinesWhitespace + b.LinesWhitespace
	r.TotalLinesProcessed = a.TotalLinesProcessed + b.TotalLinesProcessed

	r.NumInlineComments = a.NumInlineComments + b.NumInlineComments

	r.MaxCurlyBracesDepth = maxInt(a.MaxCurlyBracesDepth, b.MaxCurlyBracesDepth)
	r.MaxCurlyBracesDepthLocation = a.MaxCurlyBracesDepthLocation
	if r.MaxCurlyBracesDepth == b.MaxCurlyBracesDepth {
		r.MaxCurlyBracesDepthLocation = b.MaxCurlyBracesDepthLocation
	}

	return r
}
