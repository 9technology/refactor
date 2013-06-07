package patch

import (
	"bytes"
	"io/ioutil"
	"os"
)

func concat(first []byte, rest ...[]byte) []byte {
	for _, b := range rest {
		first = append(first, b...)
	}
	return first
}

func linesSurrounding(b []byte, start int, end int) string {
	startOfFirstLine := bytes.LastIndexAny(b[0:start], "\r\n")
	if startOfFirstLine == -1 {
		startOfFirstLine = 0
	} else {
		startOfFirstLine += 1 // Go past that newline
	}
	endOfLastLine := bytes.IndexAny(b[end:], "\r\n")
	if endOfLastLine == -1 {
		endOfLastLine = len(b)
	} else {
		endOfLastLine += end // Add the rest of the string's length
	}
	return string(b[startOfFirstLine:endOfLastLine])
}

// A Patch is an intent to replace original[start:end] with replacement
type Patch struct {
	original    []byte
	start       int
	end         int
	replacement []byte
}

func (p *Patch) Before() string {
	return linesSurrounding(p.original, p.start, p.end)
}

func (p *Patch) After() string {
	after := concat(nil, p.original[0:p.start], p.replacement, p.original[p.end:])
	return linesSurrounding(after, p.start, p.start+len(p.replacement))
}

type replacer interface {
	FindIndex([]byte) []int
	ReplaceAll([]byte, []byte) []byte
}

type Patcher struct {
	filename string
	progress int
	contents []byte
	find     replacer
	replace  []byte
}

func NewPatcher(filename string, find replacer, replace []byte) *Patcher {
	return &Patcher{filename: filename, find: find, replace: replace}
}

func (p *Patcher) Load() error {
	f, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	p.contents, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return nil
}

func (p *Patcher) Next() *Patch {
	match := p.find.FindIndex(p.contents[p.progress:])
	if match == nil {
		return nil
	}
	start := p.progress + match[0]
	end := p.progress + match[1]
	replacement := p.find.ReplaceAll(p.contents[start:end], p.replace)
	return &Patch{p.contents, start, end, replacement}
}

func (p *Patcher) Accept(patch *Patch) {
	p.contents = concat(nil, p.contents[0:patch.start], patch.replacement, p.contents[patch.end:])
	p.progress = patch.start + len(patch.replacement)
}

func (p *Patcher) Done() error {
	if p.progress == 0 {
		return nil
	}
	return ioutil.WriteFile(p.filename, p.contents, 0)
}
