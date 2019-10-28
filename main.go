package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/teejays/clog"
)

type Args struct {
	dirPath string
}

func main() {
	err := run()
	if err != nil {
		flag.PrintDefaults()
		panic(err)
	}
}

func run() error {
	// Args
	var args Args
	flag.StringVar(&args.dirPath, "dir", "", "path of the directory to analyze")
	flag.Parse()

	args.dirPath = strings.TrimSpace(args.dirPath)
	if args.dirPath == "" {
		return fmt.Errorf("directory is empty")
	}

	r, err := processDir(args.dirPath)
	if err != nil {
		return err
	}

	fmt.Printf("Results: \n%+v\n", r)

	return nil

}

type Results struct {
	LinesOfCode         int
	LinesOfErrCheck     int
	LinesOfComments     int
	LinesWhitespace     int
	TotalLinesProcessed int
}

func processDir(dirPath string) (Results, error) {
	var results Results

	// Open the directory
	clog.Debugf("Opening Dir: %s", dirPath)
	dir, err := os.Open(dirPath)
	if err != nil {
		return results, err
	}

	dInfo, err := dir.Stat()
	if err != nil {
		return results, err
	}

	if !dInfo.IsDir() {
		return results, fmt.Errorf("%s is not a directory", dirPath)
	}

	// Get names of all files
	subFiles, err := dir.Readdir(-1)
	if err != nil {
		return results, err
	}

	for _, subFile := range subFiles {
		subFilePath := joinPath(dirPath, subFile.Name())
		// If Dir
		if subFile.IsDir() {
			r, err := processDir(subFilePath)
			if err != nil {
				return results, err
			}
			results = addResults(results, r)
		}

		// If file
		r, err := processFile(subFilePath)
		if err != nil {
			return results, err
		}

		results = addResults(results, r)

	}

	return results, nil

}

func processFile(filePath string) (Results, error) {
	var r Results

	// Ignore non-Go files
	if len(filePath) < 3 || filePath[len(filePath)-3:] != ".go" {
		return r, nil
	}

	// Ignore test files
	if len(filePath) > 8 && filePath[len(filePath)-8:] == "_test.go" {
		return r, nil
	}

	// Open file
	clog.Debugf("Opening File: %s", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return r, err
	}

	stat, err := file.Stat()
	if err != nil {
		return r, err
	}

	if stat.IsDir() {
		return r, fmt.Errorf("file %s is a dir", filePath)
	}

	// Read file
	buffReader := bufio.NewReader(file)

	var lineNum int
	var inErrCheck bool
	for {
		lineNum++
		clog.Debugf("Processing Line: %d", lineNum)

		// Read the next/first line
		text, err := buffReader.ReadString('\n')
		if err != nil && err != io.EOF {
			return r, err
		}
		// End of stream
		if err == io.EOF {
			clog.Debugf("[%s] Stream EOF: %s", file.Name(), err)
			break
		}

		text = strings.TrimSpace(text)

		// Comment
		if text == "" {
			r.LinesWhitespace++

		} else if len(text) >= 2 && text[:2] == "//" {
			r.LinesOfComments++

		} else if !inErrCheck && strings.Contains(text, "if err != nil") || strings.Contains(text, "Err != nil") {
			inErrCheck = true
			r.LinesOfErrCheck++

		} else if inErrCheck && text == "}" {
			r.LinesOfErrCheck++
			inErrCheck = false

		} else {
			r.LinesOfCode++
		}

		r.TotalLinesProcessed++

	}

	err = file.Close()
	if err != nil {
		return r, err
	}
	return r, nil
}

func addResults(a, b Results) Results {
	var r Results
	r.LinesOfCode = a.LinesOfCode + b.LinesOfCode
	r.LinesOfComments = a.LinesOfComments + b.LinesOfComments
	r.LinesOfErrCheck = a.LinesOfErrCheck + b.LinesOfErrCheck
	r.LinesWhitespace = a.LinesWhitespace + b.LinesWhitespace
	r.TotalLinesProcessed = a.TotalLinesProcessed + b.TotalLinesProcessed

	return r
}

func joinPath(parts ...string) string {
	return strings.Join(parts, string(os.PathSeparator))
}
