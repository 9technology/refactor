package main

import (
	"errors"
	"fmt"
	"github.com/ninemsn/refactor/confirm"
	"github.com/ninemsn/refactor/patch"
	"github.com/vrischmann/termcolor"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func prettyPrint(filename string, patch *patch.Patch, out io.Writer) {
	fmt.Fprintf(out, "%s\n %s\n %s\n", termcolor.Colored(filename, termcolor.Cyan), termcolor.Colored("-"+patch.Before(), termcolor.Red), termcolor.Colored("+"+patch.After(), termcolor.Green))
}

func confirmPatch(filename string, p *patch.Patch, confirmation *confirm.Confirmation) bool {
	prettyPrint(filename, p, os.Stdout)
	if confirmation.Next() {
		return true
	}
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

func main() {
	if len(os.Args) <= 3 {
		println("Example: refactor .rb require import\n  Replaces 'require' with 'import' in all .rb files")
		return
	}
	suffix := os.Args[1]
	find := regexp.MustCompile(os.Args[2])
	replace := []byte(os.Args[3])

	confirmation := new(confirm.Confirmation)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		if !strings.HasSuffix(path, suffix) {
			return nil
		}
		patcher := patch.NewPatcher(path, find, replace)
		err = patcher.Load()
		if err != nil {
			return err
		}
		for p := patcher.Next(); p != nil; p = patcher.Next() {
			canProceed := confirmPatch(path, p, confirmation)
			if canProceed {
				patcher.Accept(p)
			} else {
				patcher.Done()
				return errors.New("refactor aborted by user")
			}
		}
		patcher.Done()
		return nil
	})
}
