package patch

import (
	"regexp"
	"testing"
)

func TestOuterLinesForSubstring(t *testing.T) {
	fixtures := []struct {
		Str               string
		Start             int
		ExpectedStart     int
		End               int
		ExpectedEnd       int
		ExpectedSubstring string
	}{
		{"split\nacross\ndifferent\nlines", 9, 6, 10, 12, "across"},
		{"All on the same line", 9, 0, 10, 20, "All on the same line"},
	}
	for _, fixture := range fixtures {
		lineStart, lineEnd := outerLinesForSubstring(fixture.Str, fixture.Start, fixture.End)
		if fixture.ExpectedStart != lineStart {
			t.Errorf("line containing %d should start at %d, but reported start at %d", fixture.Start, fixture.ExpectedStart, lineStart)
		}
		if fixture.ExpectedEnd != lineEnd {
			t.Errorf("line containing %d should end at %d, but reported end at %d", fixture.End, fixture.ExpectedEnd, lineEnd)
		}
		if substring := fixture.Str[lineStart:lineEnd]; fixture.ExpectedSubstring != substring {
			t.Errorf("whole line should have been %s but was %s", fixture.ExpectedSubstring, substring)
		}
	}
}

type simpleReplacer struct {
	Replace string
}

func (er simpleReplacer) ReplaceAllString(src, repl string) string {
	return er.Replace
}

func TestPatchForSubstring(t *testing.T) {
	fixtures := []struct {
		Str            string
		Start          int
		End            int
		Find           replacer
		ExpectedBefore string
		ExpectedAfter  string
		ExpectedNewStr string
	}{
		{"split\nacross\ndifferent\nlines", 7, 18, simpleReplacer{"aa"}, "across\ndifferent", "aaarent", "split\naaarent\nlines"},
		{"All on the same line", 3, 9, simpleReplacer{""}, "All on the same line", "Alle same line", "Alle same line"},
		{"Multiline replace", 3, 9, simpleReplacer{"tiple\nlines are"}, "Multiline replace", "Multiple\nlines are replace", "Multiple\nlines are replace"},
	}
	for _, fixture := range fixtures {
		newStr, before, after := patchForSubstring(fixture.Str, fixture.Start, fixture.End, fixture.Find, "")
		if newStr != fixture.ExpectedNewStr {
			t.Errorf("before should have been %s but was %s", fixture.ExpectedNewStr, newStr)
		}
		if before != fixture.ExpectedBefore {
			t.Errorf("before should have been %s but was %s", fixture.ExpectedBefore, before)
		}
		if after != fixture.ExpectedAfter {
			t.Errorf("after should have been %s but was %s", fixture.ExpectedAfter, after)
		}
	}
}

func TestPatcher(t *testing.T) {
	patches := make(chan Patch)
	proceed := make(chan bool)
	result := make(chan string)
	go Patcher("aaabbbb", regexp.MustCompile(`b`), "", patches, proceed, result)
	for patch := range patches {
		if patch.Before[0:3] != "aaa" {
			t.Errorf("patch doesn't include context")
		}
		if len(patch.Before)-len(patch.After) != 1 {
			t.Errorf("too many characters were removed")
		}
		proceed <- true
	}
	r := <-result
	if r != "aaa" {
		t.Errorf("result should be aaa but it was %s", r)
	}
}
