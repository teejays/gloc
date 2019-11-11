package main

import (
	"strings"
)

// LineResult represents the response of processing a single line
type LineResult struct {
	IsWhitespace    bool
	IsOnlyComment   bool
	IsInlineComment bool
	StartsErrCheck  bool

	NumBracesDiff int
	// HasOpenBraces  bool
	// HasCloseBraces bool

	StartsBlockComment bool
	EndsBlockComment   bool

	StartBacktickBlock bool
	EndsBacktickBlock  bool
}

// LineContext is the context that is passed to the line processor so it can process the line correctly
type LineContext struct {
	InBlockComment   bool
	InBacktickString bool
}

func processLine(text string, lineCtx LineContext) LineResult {
	var r LineResult

	// Remove surrounding whitespace
	text = strings.TrimSpace(text)

	// Line is empty
	if text == "" {
		r.IsWhitespace = true
		return r
	}

	// Line is only a comment (starts with '//')
	if len(text) >= 2 && text[:2] == "//" {
		r.IsOnlyComment = true
		return r
	}

	// vars that we need to remember as we move to the next character
	var inBlockComment, inLineComment, inDoubleQuote, inSingleQuote, inBacktick, prevWasEscapeChar bool
	var bracesDepth int
	var isInlineComment bool
	var commentIndexes = make(map[int]bool)

	// Populate any initial values
	if lineCtx.InBlockComment {
		inBlockComment = true
	}
	if lineCtx.InBacktickString {
		inBacktick = true
	}

	// Loop through the chars in the line to find braces
	for i, c := range text {

		var hasNextChar = len(text) > i+1
		var nextChar byte
		if hasNextChar {
			nextChar = text[i+1]
		}

		var hasPrevChar = i > 0
		var prevChar byte
		if hasPrevChar {
			prevChar = text[i-1]
		}

		var inString = inDoubleQuote || inSingleQuote || inBacktick
		var inComment = inBlockComment || inLineComment
		var isEscapeChar bool

		// Inline Comment Start
		if c == '/' && !inString && !inComment && hasNextChar && nextChar == '/' {
			inLineComment = true
		}
		// Block Comment Start
		var startingBlockComment = c == '/' && !inString && !inComment && hasNextChar && nextChar == '*'
		if startingBlockComment {
			inBlockComment = true
			commentIndexes[i] = true
		}
		// Block Comment End
		var endingBlockComment = c == '/' && !inString && inBlockComment && hasPrevChar && prevChar == '*'
		if endingBlockComment {
			inBlockComment = false
			commentIndexes[i] = true
		}

		// Inside a block comment
		if !startingBlockComment && !endingBlockComment && inComment {
			commentIndexes[i] = true
		}

		// if we're not ending or starting the comment, let's ignore what's in between
		if inComment {
			continue
		}

		// String Start
		if c == '"' && !inString && !inComment {
			inDoubleQuote = true
		}
		if c == '\'' && !inString && !inComment {
			inSingleQuote = true
		}
		if c == '`' && !inString && !inComment {
			inBacktick = true
		}

		// String End
		if c == '"' && inString && inDoubleQuote && !inComment && !prevWasEscapeChar {
			inDoubleQuote = false
		}
		if c == '\'' && inString && inSingleQuote && !inComment && !prevWasEscapeChar {
			inSingleQuote = false
		}
		if c == '`' && inString && inBacktick && !inComment {
			inBacktick = false
		}

		// Escape Char
		if c == '\\' && !inComment && (inDoubleQuote || inSingleQuote) && !prevWasEscapeChar {
			isEscapeChar = true
		}

		// Curly Brackets
		if c == '{' && !inString && !inComment {
			bracesDepth++
		}
		if c == '}' && !inString && !inComment {
			bracesDepth--
		}

		prevWasEscapeChar = isEscapeChar
		if inLineComment {
			isInlineComment = true
			text = text[:i] // remaining part is only a comment, so skip it
			break           // no need to process remaining comment chars
		}
	}

	// How does the text look without the block comments
	if len(commentIndexes) > 0 {
		var copyText []rune
		for i, r := range text {
			if !commentIndexes[i] {
				copyText = append(copyText, r)
			}
		}
		// Remove surrounding whitespace
		cleanText := strings.TrimSpace(string(copyText))
		if cleanText == "" {
			r.IsOnlyComment = true
		}
		if cleanText != "" {
			r.IsInlineComment = true
		}
	}
	if isInlineComment {
		r.IsInlineComment = true
	}

	if inBlockComment {
		r.StartsBlockComment = true
	}
	if lineCtx.InBlockComment && !inBlockComment {
		r.EndsBlockComment = true
	}

	if inBacktick {
		r.StartBacktickBlock = true
	}
	if lineCtx.InBacktickString && !inBacktick {
		r.EndsBacktickBlock = true
	}

	// We're starting an err check
	if strings.Contains(text, "if err != nil") || strings.Contains(text, "Err != nil") {
		r.StartsErrCheck = true
	}

	r.NumBracesDiff = bracesDepth

	return r
}
