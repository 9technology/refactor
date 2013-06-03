package patch

import (
	"fmt"
	"strings"
)

func outerLinesForSubstring(str string, start int, end int) (startOfFirstLine int, endOfLastLine int) {
	startOfFirstLine = strings.LastIndex(str[0:start], "\n")
	if startOfFirstLine == -1 {
		startOfFirstLine = 0
	} else {
		startOfFirstLine += 1 // Go past that newline
	}
	endOfLastLine = strings.Index(str[end:], "\n")
	if endOfLastLine == -1 {
		endOfLastLine = len(str)
	} else {
		endOfLastLine += end // Add the rest of the string's length
	}
	return
}

type replacer interface {
	ReplaceAllString(src, repl string) string
}

// Returns a before/after view of the surrounding lines when find/replace is done
func patchForSubstring(str string, start int, end int, find replacer, replace string) (newStr string, before string, after string) {
	replaced := find.ReplaceAllString(str[start:end], replace)
	originalStart, originalEnd := outerLinesForSubstring(str, start, end)
	before = str[originalStart:originalEnd]
	newStr = str[0:start] + replaced + str[end:]
	newStart, newEnd := outerLinesForSubstring(newStr, start, start+len(replaced))
	after = newStr[newStart:newEnd]
	return
}

type Patch struct {
	Filename string
	Before   string
	After    string
}

func (p Patch) String() string {
	return fmt.Sprintf("%s:\n\t%s => %s\n", p.Filename, p.Before, p.After)
}

type finder interface {
	replacer
	FindStringIndex(str string) []int
}

func Patcher(str string, find finder, replace string, patches chan<- Patch, proceed <-chan bool, result chan<- string) {
	defer close(result)
	var before, after string
	var progress int
	match := find.FindStringIndex(str[progress:])
	if match == nil {
		close(patches)
		return
	}
	for match != nil {
		str, before, after = patchForSubstring(str, progress+match[0], progress+match[1], find, replace)
		progress = match[1] - len(before) + len(after)
		patches <- Patch{Before: before, After: after}
		canProceed := <-proceed
		if !canProceed {
			break
		}
		match = find.FindStringIndex(str[progress:])
	}
	close(patches)
	result <- str
}
