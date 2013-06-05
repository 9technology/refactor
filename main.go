package main

import (
	"./confirm"
	"fmt"
	"github.com/pranavraja/refactor/patch"
	"github.com/vrischmann/termcolor"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

func walk(root string, suffix string, filePaths chan<- string) {
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
			go walk(fullPath, suffix, nestedFilePaths)
			for f := range nestedFilePaths {
				filePaths <- f
			}
		} else {
			if strings.HasSuffix(fullPath, suffix) {
				filePaths <- fullPath
			}
		}
	}
}

func patchAll(filenames <-chan string, find *regexp.Regexp, replace string, patches chan<- patch.Patch, proceed <-chan bool) {
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
		filePatches := make(chan patch.Patch)
		patcherCanProceed := make(chan bool)
		fileResult := make(chan string)
		go patch.Patcher(string(contents), find, replace, filePatches, patcherCanProceed, fileResult)
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

func prettyPrint(patch patch.Patch, out io.Writer) {
	fmt.Fprintf(out, "%s\n %s\n %s\n", termcolor.Colored(patch.Filename, termcolor.Cyan), termcolor.Colored("-"+patch.Before, termcolor.Red), termcolor.Colored("+"+patch.After, termcolor.Green))
}

func main() {
	if len(os.Args) <= 3 {
		println("Example: refactor .rb require import\n  Replaces 'require' with 'import' in all .rb files")
		return
	}
	suffix := os.Args[1]
	find := regexp.MustCompile(os.Args[2])
	replace := os.Args[3]
	patches := make(chan patch.Patch)
	proceed := make(chan bool)

	paths := make(chan string)
	go walk(".", suffix, paths)
	go patchAll(paths, find, replace, patches, proceed)
	var confirmation confirm.Confirmation
	for p := range patches {
		prettyPrint(p, os.Stdout)
		if confirmation.Next() {
			proceed <- true
		} else {
			fmt.Printf("Continue? ([a]ll/[y]es/[n]o (default no): ")
			var input string
			_, err := fmt.Scanf("%s", &input)
			if err != nil {
				return
			}
			switch input {
			case "a":
				confirmation.ConfirmAll()
			case "y":
				confirmation.ConfirmOnce()
			default:
				return
			}
			proceed <- confirmation.Next()
		}
	}
}
