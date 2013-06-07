package patch

import (
	"regexp"
	"testing"
)

var singleLinePatch []byte = []byte("single line")
var multilinePatch []byte = []byte("split\nacross\ndifferent\nlines")

func TestPatch_Before(t *testing.T) {
	var patch Patch
	patch = Patch{
		original: singleLinePatch,
		start:    4,
		end:      5,
	}
	if before := patch.Before(); before != "single line" {
		t.Errorf("patch.Before() should be 'single line', but was '%s'", before)
	}
	patch = Patch{
		original: multilinePatch,
		start:    8,
		end:      12,
	}
	if before := patch.Before(); before != "across" {
		t.Errorf("patch.Before() should be 'across', but was '%s'", before)
	}
	patch = Patch{
		original: multilinePatch,
		start:    8,
		end:      16,
	}
	if before := patch.Before(); before != "across\ndifferent" {
		t.Errorf("patch.Before() should be 'across\ndifferent', but was '%s'", before)
	}
}

func TestPatch_After(t *testing.T) {
	var patch Patch
	patch = Patch{
		original:    singleLinePatch,
		start:       4,
		end:         5,
		replacement: []byte("replacement"),
	}
	if after := patch.After(); after != "singreplacemente line" {
		t.Errorf("patch.After() should be 'singreplacemente line', but was '%s'", after)
	}
	patch = Patch{
		original:    multilinePatch,
		start:       8,
		end:         12,
		replacement: []byte("replacement"),
	}
	if after := patch.After(); after != "acreplacement" {
		t.Errorf("patch.After() should be 'acreplacement', but was '%s'", after)
	}
	patch = Patch{
		original:    multilinePatch,
		start:       8,
		end:         16,
		replacement: []byte("replacement"),
	}
	if after := patch.After(); after != "acreplacementferent" {
		t.Errorf("patch.After() should be 'acreplacementferent', but was '%s'", after)
	}
	patch = Patch{
		original:    multilinePatch,
		start:       8,
		end:         16,
		replacement: []byte("multiline\nreplace"),
	}
	if after := patch.After(); after != "acmultiline\nreplaceferent" {
		t.Errorf("patch.After() should be 'acmultiline\nreplaceferent', but was '%s'", after)
	}
}

func TestPatcher_Next(t *testing.T) {
	patcher := Patcher{
		find:     regexp.MustCompile(`ab`),
		replace:  []byte("cd"),
		contents: []byte("gabcdefgabcdef"),
	}
	patch := patcher.Next()
	if string(patch.original) != "gabcdefgabcdef" {
		t.Errorf("string was not passed to *Patch")
	}
	if patch.start != 1 {
		t.Errorf("patch start %d is incorrect", patch.start)
	}
	if patch.end != 3 {
		t.Errorf("patch end %d is incorrect", patch.end)
	}
	if string(patch.replacement) != "cd" {
		t.Errorf("incorrect patch replacement")
	}
}

func TestPatcher_Next_NotFound(t *testing.T) {
	patcher := Patcher{
		find:     regexp.MustCompile(`[0-9]`),
		contents: []byte("gabcdefgabcdef"),
	}
	if next := patcher.Next(); next != nil {
		t.Errorf("didn't return nil on not found: %v", next)
	}
}

func TestPatcher_Next_CaptureGroups(t *testing.T) {
	patcher := Patcher{
		find:     regexp.MustCompile(`a(\w\w)`),
		replace:  []byte("c$1"),
		contents: []byte("gabcdefgabcdef"),
	}
	patch := patcher.Next()
	if string(patch.original) != "gabcdefgabcdef" {
		t.Errorf("string was not passed to *Patch")
	}
	if patch.start != 1 {
		t.Errorf("patch start %d is incorrect", patch.start)
	}
	if patch.end != 4 {
		t.Errorf("patch end %d is incorrect", patch.end)
	}
	if string(patch.replacement) != "cbc" {
		t.Errorf("incorrect patch replacement")
	}
}

func TestPatcher_Accept(t *testing.T) {
	patcher := Patcher{
		contents: []byte("gabcdefgabcdef"),
	}
	patch := &Patch{
		start:       1,
		end:         10,
		replacement: []byte("asd"),
	}
	patcher.Accept(patch)
	if str := string(patcher.contents); str != "gasdcdef" {
		t.Errorf("patch incorrectly applied: %s", str)
	}
	if patcher.progress != 4 {
		t.Errorf("incorrect progress through file: %d", patcher.progress)
	}
}

func TestPatcher_Done_NoChanges(t *testing.T) {
	patcher := Patcher{}
	if err := patcher.Done(); err != nil {
		t.Errorf("error on Done(): %v", err)
	}
}
