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
	rootPath        string
	excludeDirs     string
	excludeFiles    string
	ignoreTestFiles bool
}

func main() {
	err := run()
	if err != nil {
		flag.PrintDefaults()
		panic(err)
	}
}

func run() error {
	// Parse Args
	var args Args
	flag.StringVar(&args.rootPath, "root", "", "path of the directory to analyze")
	flag.BoolVar(&args.ignoreTestFiles, "ignore-test-files", true, "should be ignore test files (default to true)")
	flag.StringVar(&args.excludeDirs, "exclude-dirs", "", "directories to be excluded (comma separated)")
	flag.StringVar(&args.excludeFiles, "exclude-files", "", "files to be excluded (comma separated)")
	flag.Parse()

	args.rootPath = strings.TrimSpace(args.rootPath)
	if args.rootPath == "" {
		return fmt.Errorf("directory is empty")
	}

	// Process the root project directory
	config := fileConfig{
		ignoreTestFiles: args.ignoreTestFiles,
		excludeDirs:     strings.Split(args.excludeDirs, ","),
		excludeFiles:    strings.Split(args.excludeFiles, ","),
	}
	r, err := processDir(args.rootPath, config)
	if err != nil {
		return err
	}

	fmt.Printf("Results: \n%+v\n", r)

	return nil

}

type fileConfig struct {
	excludeDirs     []string
	excludeFiles    []string
	ignoreTestFiles bool
}

func processDir(dirPath string, config fileConfig) (Results, error) {
	var results Results

	// Excluded dirs
	if sliceContainsString(config.excludeDirs, dirPath) {
		return results, nil
	}

	// Open the directory
	clog.Debugf("Opening Dir: %s", dirPath)
	dir, err := os.Open(dirPath)
	if err != nil {
		return results, err
	}

	// Find whether the file is a dir or not.
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

		// If Dir
		if subFile.IsDir() {
			r, err := processDir(joinPath(dirPath, subFile.Name()), config)
			if err != nil {
				return results, err
			}
			results = addResults(results, r)
		}

		// If file
		r, err := processFile(dirPath, subFile.Name(), config)
		if err != nil {
			return results, err
		}

		results = addResults(results, r)

	}

	return results, nil

}

func processFile(dirPath, fileName string, config fileConfig) (Results, error) {
	var r Results

	// Open file
	filePath := joinPath(dirPath, fileName)
	clog.Debugf("Opening File: %s", filePath)

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

	r.NumOfFiles++
	r.MaxCurlyBracesDepthLocation.File = filePath

	// Read file
	buffReader := bufio.NewReader(file)

	// Scope level (count number of '{' that we are in)
	var errCheckDepth int
	var curlyBracesDepth int

	// type BraceType string
	// var curlyBracesStack = NewStack()

	var lineNum int

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

		// Line is empty
		if text == "" {
			r.LinesWhitespace++
			continue

		}

		// Line is only a comment (starts with '//')
		if len(text) >= 2 && text[:2] == "//" {
			r.LinesOfComments++
			continue

		}

		// Loop through the chars in the line
		var hasComment bool
		for i, c := range text {

			if c == '{' {
				curlyBracesDepth++
			}
			if c == '}' {
				curlyBracesDepth--
			}

			if c == '/' && len(text) > i+1 && text[i+1] == '/' {
				text = text[:i] // this is the line text without the comment part
				hasComment = true
				break // no need to process comments
			}
		}

		if hasComment {
			r.NumInlineComments++
		}

		if curlyBracesDepth > r.MaxCurlyBracesDepth {
			r.MaxCurlyBracesDepth = curlyBracesDepth
			r.MaxCurlyBracesDepthLocation.Line = lineNum
		}

		var isErrCheckLine bool

		// We're starting an err check
		if strings.Contains(text, "if err != nil") || strings.Contains(text, "Err != nil") {
			isErrCheckLine = true
			errCheckDepth++
		}

		if errCheckDepth > 0 && text != "}" {
			isErrCheckLine = true
		}

		// Closing an error check
		if errCheckDepth > 0 && text == "}" {
			isErrCheckLine = true
			errCheckDepth--
		}

		if isErrCheckLine {
			r.LinesOfErrCheck++
		}

		if !isErrCheckLine {
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

func joinPath(parts ...string) string {
	return strings.Join(parts, string(os.PathSeparator))
}

func sliceContainsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func maxInt(arr ...int) int {
	if len(arr) < 1 {
		panic("maxInt called with no ints")
	}

	var max = arr[0]
	for _, n := range arr {
		if n > max {
			max = n
		}
	}
	return max
}
