package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandAdd,
	commandList,
	commandEdit,
	commandConfig,
}

var commandAdd = cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Flags:   Flags,
	Usage:   "Add bookmark",
	Description: `
`,
	Action: doAdd,
}

var commandList = cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Show bookmark list",
	Flags:   Flags,
	Description: `
`,
	Action: doList,
}

var commandEdit = cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit bookmark list",
	Flags:   Flags,
	Description: `
`,
	Action: doEdit,
}

var commandConfig = cli.Command{
	Name:  "config",
	Usage: "Set configure",
	Description: `gookmark config ui.editor=vim
   gookmark config core.linefeed=dos
`,
	Action: doConfig,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doAdd(c *cli.Context) {
	group := strings.Replace(c.String("group"), "\"", "", -1)

	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var linefeed string
	if config.Core.Linefeed == "unix" {
		linefeed = "\n"
	} else if config.Core.Linefeed == "dos" {
		linefeed = "\r\n"
	} else {
		linefeed = "\n"
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var item string
	if len(c.Args()) == 0 {
		item = wd
	} else {
		for _, arg := range c.Args() {
			runes := []rune(arg)
			if runes[0] == '/' || runes[0] == '\\' {
				if runtime.GOOS == "windows" {
					if runes[0] == runes[1] {
						item = arg
					} else {
						drive := strings.Split(wd, string(os.PathSeparator))[0]
						item = filepath.Join(drive, arg)
					}
				} else {
					item = arg
				}
			} else if filepath.IsAbs(arg) {
				item = arg
			} else {
				item = filepath.Join(wd, arg)
			}
		}
	}
	item = filepath.Clean(item)

	bookmarks, err := GetBookmarks(group)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var uniqBookmarks []string
	for _, bookmark := range bookmarks {
		if bookmark != item {
			uniqBookmarks = append(uniqBookmarks, bookmark)
		}
	}
	uniqBookmarks = append([]string{item}, uniqBookmarks...)

	path, err := GetBookmarkFilePath(group)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	os.Mkdir(filepath.Dir(path), 0777)
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer fp.Close()

	fmt.Fprint(fp, strings.Join(uniqBookmarks, linefeed))
}

func doList(c *cli.Context) {
	group := strings.Replace(c.String("group"), "\"", "", -1)

	bookmarks, err := GetBookmarks(group)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, bookmark := range bookmarks {
		fmt.Println(bookmark)
	}
}

func doEdit(c *cli.Context) {
	group := strings.Replace(c.String("group"), "\"", "", -1)

	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	path, err := GetBookmarkFilePath(group)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	cmd := exec.Command(config.Ui.Editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func doConfig(c *cli.Context) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if len(c.Args()) == 0 {
		PrintCurrentConfig(&config)
	} else {
		for _, arg := range c.Args() {
			err := WriteNewConfig(&config, arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func GetHomePath() (string, error) {
	var home string
	var err error
	home = os.Getenv("HOME")
	if len(home) == 0 {
		home = os.Getenv("USERPROFILE")
	}
	if len(home) == 0 {
		err = errors.New("Not fond your HOME path!")
	}
	return home, err
}

func GetBookmarkFilePath(group string) (string, error) {
	var path string
	home, err := GetHomePath()
	if err == nil {
		appData := filepath.Join(home, ".gookmark")
		path = filepath.Join(appData, group+".txt")
	}
	return path, err
}

func GetBookmarks(group string) ([]string, error) {
	var bookmarks []string
	path, err := GetBookmarkFilePath(group)
	if err == nil {
		fp, err := os.Open(path)
		if err == nil {
			defer fp.Close()
			scanner := bufio.NewScanner(fp)
			for scanner.Scan() {
				bookmarks = append(bookmarks, strings.TrimRight(strings.TrimRight(scanner.Text(), "\n"), "\r"))
			}
		}
	}
	return bookmarks, err
}

func PrintCurrentConfig(config *Config) {
	fmt.Printf("ui.editor=%q\n", config.Ui.Editor)
	fmt.Printf("core.linefeed=%q\n", config.Core.Linefeed)
	for k, v := range config.Alias {
		fmt.Printf("alias.%s=%q\n", k, v)
	}
}
