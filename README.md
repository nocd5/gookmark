# gookmark
Bookmark for CLI

## Description

## Usage

```bash
# current directory : /foo/bar/hoge/fuga
$ gookmark add
$ gookmark add ..
$ gookmark list
#   /foo/bar/hoge/fuga
#   /foo/bar/hoge
$ gookmark edit
```

### bookmark group
`gookmark <subcommand> --group <group name> <file or directory>`  
default &lt;group name&gt; is **bookmark**
```bash
# current directory : /foo/bar/hoge/fuga
$ gookmark add --group files foo.txt
$ gookmark add  --group files bar.txt
$ gookmark list --group files
#   /foo/bar/hoge/fuga/foo.txt
#   /foo/bar/hoge/fuga/bar.txt
$ gookmark edit --group files
```

### subcommand alias
support alias for subcommand
```bash
# set alias
$ gookmark config alias.af="add --group f"
$ gookmark config alias.al="list --group f"
$ gookmark config alias.ae="edit --group f"
# use alias
$ gookmark af foo.txt
$ gookmark af files bar.txt
$ gookmark al files
#   /foo/bar/hoge/fuga/foo.txt
#   /foo/bar/hoge/fuga/bar.txt
$ gookmark ae files
```
**alias** should be set `subcommand` + `bookmark group option`

## Install

To install, use `go get`:

```bash
$ go get -d github.com/nocd5/gookmark
$ cd $GOPATH/src/github.com/nocd5/gookmark
$ go install
$ gookmark config ui.editor=vim
$ gookmark config core.linefeed=unix
```
`gookmark config ui.editor=your_editor`  
`gookmark config core.linefeed=[unix|dos]`

## Contribution

1. Fork ([https://github.com/nocd5/gookmark/fork](https://github.com/nocd5/gookmark/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `go fmt`
1. Create a new Pull Request

## Author

[nocd5](https://github.com/nocd5)
