package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/teejays/clog"
)

// Args can be passed as the command line arguments, and control the program
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
