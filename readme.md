# GLOC
Go Lines of Code (GLOC) is a simple executable tool that runs basic code analysis on a Go project. 

### Why did I write this?

I once did a small take-home assignment for which I ended up spending more time and writing more code than I had planned. In the end, I wanted to some kind of analysis around if I have written too much (or too less) code. 

There are tools out there that measure lines of codes in programs but I couldn't find any that was Go specific. As we know, _n_ lines of code in Go don't necessarily correspond to the same amount of code as _n_ lines of let's say Java. This is because of the idiomatic error checking in Go that takes a few lines after most function calls, and because of test files that are part of the same package. 

### What does it do?
It can tell you the following about your code:

- Number of total lines
- Number of lines relevant code (excluding error checking)
- Number of lines of code that corresponds to error checking
- Number of whitespace lines
- Number of lines that are pure comments
- Number of lines that have inline comments
- Maximum scope depth (i.e. how many nested levels of curly braces do we go) and where

## Getting Started

### Prerequisites
- You should have go installed

### Setting Up

1) Install Gloc by running:

``` go get -u github.com/teejays/gloc ```

2) 
    1) If you have Gobin setup and added to your path, running the above command would automatically allow you to run "`gloc help`".  If successful, ignore the next step.
    
    2) If not successful, you could manually build Gloc by navigating to the directory in which Gloc code is checked out, usually `$GOPATH/src/github.com/teejays/gloc`. Once in  that directory, you can build the binary by running `make build`, and then running "`./bin/gloc help`".

### Usage
Once verified that Gloc is installed, run it like this:

```gloc --root=<dir with some Go code> --ignore-test-files=<true/false> --exclude-dirs=<a,b> --exclude-files=<a.go,b.go>```

(replace `gloc` with  `./bin/gloc` if you built the binary yourself using Step 2.2 above)

**Sample Output**:

```
{
    NumOfFiles:10 
    LinesOfCode:1048
    LinesOfErrCheck:118
    LinesOfComments:188
    LinesWhitespace:320
    TotalLinesProcessed:1166
    NumInlineComments:18
    MaxCurlyBracesDepth:5
    MaxCurlyBracesDepthLocation: {
        File:../logdog/main.go 
        Line:191
    }
}
```

## Issues & Bugs

Please feel free to open Github Issues or make Pull Requests if you find any bug or need to add features.