package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func processFile(dirPath, fileName string, config fileConfig) (Results, error) {

	var r Results

	if !shouldIncludeFile(dirPath, fileName, config) {
		return r, nil
	}

	// Construct file path
	filePath := joinPath(dirPath, fileName)

	// Open File
	file, err := os.Open(filePath)
	if err != nil {
		return r, err
	}

	// Get Buf.Reader
	bufReader, err := getBufReader(file)
	if err != nil {
		return r, err
	}

	r, err = processBufReader(bufReader)
	if err != nil {
		return r, fmt.Errorf("file %s: %s", file.Name(), err)
	}

	// Close the file
	err = file.Close()
	if err != nil {
		return r, err
	}

	r.MaxCurlyBracesDepthLocation.File = filePath

	return r, nil
}

func shouldIncludeFile(dirPath, fileName string, config fileConfig) bool {
	// Ignore non-Go files
	if len(fileName) < 3 || fileName[len(fileName)-3:] != ".go" {
		return false
	}

	// Ignore test files
	if config.ignoreTestFiles && len(fileName) > 8 && fileName[len(fileName)-8:] == "_test.go" {
		return false
	}

	// Excluded files
	filePath := joinPath(dirPath, fileName)
	if sliceContainsString(config.excludeFiles, filePath) {
		return false
	}

	return true
}

func getBufReader(file *os.File) (*bufio.Reader, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("file %s: is a dir", file.Name())
	}

	// Read file
	buffReader := bufio.NewReader(file)

	return buffReader, nil

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

		// Read the next/first line
		text, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return r, err
		}
		// End of stream
		if err == io.EOF {
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

	}

	return r, nil
}
