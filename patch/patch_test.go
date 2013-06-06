package patch

import (
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
		lines := string(linesSurrounding([]byte(fixture.Str), fixture.Start, fixture.End))
		if fixture.ExpectedSubstring != lines {
			t.Errorf("whole line should have been %s but was %s", fixture.ExpectedSubstring, lines)
		}
	}
}
