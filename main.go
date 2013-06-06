package main

import (
	"fmt"
	"github.com/pranavraja/refactor/confirm"
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

func prettyPrint(filename string, patch *patch.Patch, out io.Writer) {
	fmt.Fprintf(out, "%s\n %s\n %s\n", termcolor.Colored(filename, termcolor.Cyan), termcolor.Colored("-"+patch.Before(), termcolor.Red), termcolor.Colored("+"+patch.After(), termcolor.Green))
}

func confirmPatch(filename string, p *patch.Patch, confirmation *confirm.Confirmation) bool {
	prettyPrint(filename, p, os.Stdout)
	if confirmation.Next() {
		return true
	} else {
		fmt.Printf("Continue? ([a]ll/[y]es/[n]o (default no): ")
		var input string
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			return false
		}
		switch input {
		case "a":
			confirmation.ConfirmAll()
		case "y":
			confirmation.ConfirmOnce()
		default:
			return false
		}
		return confirmation.Next()
	}
}

func main() {
	if len(os.Args) <= 3 {
		println("Example: refactor .rb require import\n  Replaces 'require' with 'import' in all .rb files")
		return
	}
	suffix := os.Args[1]
	find := regexp.MustCompile(os.Args[2])
	replace := []byte(os.Args[3])

	paths := make(chan string)
	go walk(".", suffix, paths)
	confirmation := new(confirm.Confirmation)
	for file := range paths {
		patcher := patch.NewPatcher(file, find, replace)
		err := patcher.Load()
		if err != nil {
			println(err)
			return
		}
		for p := patcher.Next(); p != nil; p = patcher.Next() {
			canProceed := confirmPatch(file, p, confirmation)
			if canProceed {
				patcher.Accept(p)
			} else {
				patcher.Done()
				return
			}
		}
	}
}
