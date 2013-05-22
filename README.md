
A tool inspired by [codemod](https://github.com/facebook/codemod) to automate refactorings across your codebase. Useful when working without an IDE or with XCode. Skips hidden files and files without an extension.

`go get github.com/pranavraja/refactor`

## Usage

`refactor .ext regexp replacement`

Apply `s/regexp/replacement/g` across all files in the current directory tree with extension `.ext`

`regexp` is parsed by the [Go regexp package](http://golang.org/pkg/regexp/), so you can probably use any of the [re2](https://code.google.com/p/re2/wiki/Syntax) syntax. For example, you can use `$1`-style placeholders to refer to capture groups.

## Why not codemod?

Codemod asks me for confirmation for every change, which I find tiring. Usually once I see the first change I either want to quit or apply the rest of the changes.

## Why not use the power/elegance/flexibility/chainsaw of unix?

For example:

- `find . -name '$1' | xargs sed -i "" 's/$2/$3/g'` 

This is great, but I still want to see the first change to determine whether to back out. This doesn't give me confirmation.

- `grep -ERli --include "$1" "$2" . | xargs -o vim -c "argdo %s/$2/$3/gce | update" -c 'q'`

This one allows me to 'apply all' or 'skip all' on a per-file basis, not across the project. Also with a large number of files it tends to make vim freak out and write swap files all over the place.

## Running the tests

`go test ./src/refactor`

## TODO

- Prettier display of patch hunks

