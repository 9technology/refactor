
A tool to automate regex-based refactorings across your codebase. Similar to [codemod](https://github.com/facebook/codemod), but slightly more optimistic.

`go get github.com/ninemsn/refactor`

## Usage

Assuming that `$GOPATH/bin` is in your `$PATH`:

`refactor .ext 'regexp' 'replacement'`

Apply `s/regexp/replacement/g` across all files in the current directory tree with extension `.ext`. Skip hidden files. Output each change applied as a color-coded patch, and allow the user to confirm all, confirm one, or abort.

`regexp` is parsed by the [Go regexp package](http://golang.org/pkg/regexp/), so you can probably use any of the [re2](https://code.google.com/p/re2/wiki/Syntax) syntax. For example, you can use `$1`-style placeholders to refer to capture groups, as long as you single-quote the replacement string so that your shell doesn't interpolate `$1` as something else. 

## Why not use the power/elegance/flexibility/chainsaw of unix?

For example:

- `find . -name '$1' | xargs sed -i "" 's/$2/$3/g'` 

This is great, but I still want to see the first change to determine whether to back out. This doesn't give me confirmation.

- `grep -ERli --include "$1" "$2" . | xargs -o vim -c "argdo %s/$2/$3/gce | update" -c 'q'`

This one allows me to 'apply all' or 'skip all' on a per-file basis, not across the project. Also with a large number of files it tends to make vim freak out and write swap files all over the place.

## Running the tests

`go test ./patch`
`go test ./confirm`

## TODO

- Prettier display of patch hunks

