package main

import (
	"github.com/pranavraja/refactor/src/refactor"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

func walk(root string, filePaths chan<- string) {
	defer close(filePaths)
	entries, err := ioutil.ReadDir(root)
	if err != nil {
		return
	}
	for _, file := range entries {
		if strings.HasPrefix(file.Name(), ".") {
			continue // Ignore hidden files
		}
		fullPath := path.Join(root, file.Name())
		if file.IsDir() {
			nestedFilePaths := make(chan string)
			go walk(fullPath, nestedFilePaths)
			for f := range nestedFilePaths {
				filePaths <- f
			}
		} else {
			if !strings.Contains(file.Name(), ".") {
				continue // Ignore binary files
			}
			filePaths <- fullPath
		}
	}
}

func patchAll(filenames <-chan string, find *regexp.Regexp, replace string, patches chan<- refactor.Patch, proceed <-chan bool) {
	defer close(patches)
	for filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			continue
		}
		contents, err := ioutil.ReadAll(f)
		f.Close()
		if err != nil {
			continue
		}
		filePatches := make(chan refactor.Patch)
		patcherCanProceed := make(chan bool)
		fileResult := make(chan string)
		go refactor.Patcher(string(contents), find, replace, filePatches, patcherCanProceed, fileResult)
		for patch := range filePatches {
			patch.Filename = filename
			patches <- patch
			patcherCanProceed <- <-proceed
		}
		r, changed := <-fileResult
		if changed {
			ioutil.WriteFile(filename, []byte(r), 0) // Assume file already exists
		}
	}
}

func main() {
	if len(os.Args) <= 2 {
		println("Usage: refactor find replace")
		return
	}
	find := regexp.MustCompile(os.Args[1])
	replace := os.Args[2]
	patches := make(chan refactor.Patch)
	proceed := make(chan bool)

	paths := make(chan string)
	go walk(".", paths)
	go patchAll(paths, find, replace, patches, proceed)
	var canProceed bool
	for p := range patches {
		fmt.Printf("%v\n", p)
		if !canProceed {
			fmt.Printf("Continue? (y/n[default]): ")
			var input rune
			_, err := fmt.Scanf("%c", &input)
			if input == 'y' {
				canProceed = true
			}
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
			if !canProceed {
				return
			}
		}
		proceed <- true
	}
}
