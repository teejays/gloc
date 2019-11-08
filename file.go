package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func processFile(dirPath, fileName string, config fileConfig) (Results, error) {

	var r Results

	// Construct file path
	filePath := joinPath(dirPath, fileName)

	// Ignore non-Go files
	if len(fileName) < 3 || fileName[len(fileName)-3:] != ".go" {
		return r, nil
	}

	// Ignore test files
	if config.ignoreTestFiles && len(fileName) > 8 && fileName[len(fileName)-8:] == "_test.go" {
		return r, nil
	}

	// Excluded files
	if sliceContainsString(config.excludeFiles, filePath) {
		return r, nil
	}

	// Open File
	file, err := os.Open(filePath)
	if err != nil {
		return r, err
	}

	r, err = processFileReader(file)
	if err != nil {
		return r, err
	}

	r.MaxCurlyBracesDepthLocation.File = filePath

	return r, nil
}

func processFileReader(file *os.File) (Results, error) {
	var r Results

	stat, err := file.Stat()
	if err != nil {
		return r, err
	}

	if stat.IsDir() {
		return r, fmt.Errorf("file %s: is a dir", file.Name())
	}

	// Read file
	buffReader := bufio.NewReader(file)

	r, err = processBufReader(buffReader)
	if err != nil {
		return r, fmt.Errorf("file %s: %s", file.Name(), err)
	}

	err = file.Close()
	if err != nil {
		return r, err
	}

	return r, nil

}

func processBufReader(reader *bufio.Reader) (Results, error) {
	var r Results

	r.NumOfFiles = 1

	var lineNum int
	var bracesDepth int

	var errCheckPoint int // the bracesDepth at which error check starts

	var lineCtx LineContext

	for {

		lineNum++

		var isLastLine bool
		// Read the next/first line
		text, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return r, err
		}
		// End of stream
		if err == io.EOF {
			isLastLine = true
			break
		}

		r.TotalLinesProcessed++

		lr := processLine(text, lineCtx)

		// Line is empty
		if lr.IsWhitespace {
			r.LinesWhitespace++
			continue
		}

		// Line is only a comment
		if lr.IsOnlyComment && lr.IsInlineComment {
			return r, fmt.Errorf("line %d: IsOnlyComment and IsInlineComment both as true", lineNum)
		}

		if lr.IsOnlyComment {
			r.LinesOfComments++
			continue
		}

		if lr.IsInlineComment {
			r.NumInlineComments++
		}

		// Block Comments
		if lr.StartsBlockComment && lr.EndsBlockComment {
			return r, fmt.Errorf("line %d: StartsBlockComment and EndsBlockComment both as true", lineNum)
		}
		if lr.StartsBlockComment {
			lineCtx.InBlockComment = true
		}
		if lr.EndsBlockComment {
			lineCtx.InBlockComment = false
		}

		// Backtick Multiline Strings
		if lr.StartBacktickBlock && lr.EndsBacktickBlock {
			return r, fmt.Errorf("line %d: returned StartBacktickBlock and EndsBacktickBlock both as true", lineNum)
		}
		if lr.StartBacktickBlock {
			lineCtx.InBacktickString = true
		}
		if lr.EndsBacktickBlock {
			lineCtx.InBacktickString = false
		}

		// Handle Braces Depth
		bracesDepth = bracesDepth + lr.NumBracesDiff

		if bracesDepth > r.MaxCurlyBracesDepth {
			r.MaxCurlyBracesDepth = bracesDepth
			r.MaxCurlyBracesDepthLocation.Line = lineNum
		}

		if errCheckPoint == 0 && lr.StartsErrCheck {
			errCheckPoint = bracesDepth - 1
			r.LinesOfErrCheck++

		} else if errCheckPoint > 0 && lr.NumBracesDiff < 0 && bracesDepth == errCheckPoint {
			errCheckPoint = 0
			r.LinesOfErrCheck++

		} else if errCheckPoint > 0 {
			r.LinesOfErrCheck++
		}

		if errCheckPoint == 0 {
			r.LinesOfCode++
		}

		if isLastLine {
			break
		}

	}

	return r, nil
}
