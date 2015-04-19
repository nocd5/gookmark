package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
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
	Usage:   "Add bookmark",
	Description: `
`,
	Action: doAdd,
}

var commandList = cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Show bookmark list",
	Description: `
`,
	Action: doList,
}

var commandEdit = cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit bookmark list",
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

	bookmarks, err := GetBookmarks()
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

	path, err := GetBookmarkFilePath()
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
	bookmarks, err := GetBookmarks()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, bookmark := range bookmarks {
		fmt.Println(bookmark)
	}
}

func doEdit(c *cli.Context) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	path, err := GetBookmarkFilePath()
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
		ShowCurrentConfig(&config)
	} else {
		for _, arg := range c.Args() {
			SetNewConfig(&config, arg)
		}
	}
}

func GetHomePath() (string, error) {
	var home string = ""
	var err error = nil
	home = os.Getenv("HOME")
	if len(home) == 0 {
		home = os.Getenv("USERPROFILE")
	}
	if len(home) == 0 {
		err = errors.New("Not fond your HOME path!")
	}
	return home, err
}

func GetBookmarkFilePath() (string, error) {
	var path string = ""
	home, err := GetHomePath()
	if err == nil {
		appData := filepath.Join(home, ".gookmark")
		path = filepath.Join(appData, "bookmark.txt")
	}
	return path, err
}

func GetBookmarks() ([]string, error) {
	var bookmarks []string = nil
	path, err := GetBookmarkFilePath()
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

type Config struct {
	Ui   UiSection   `toml:"ui"`
	Core CoreSection `toml:"core"`
}

type UiSection struct {
	Editor string `toml:"editor"`
}

type CoreSection struct {
	Linefeed string `toml:"linefeed"`
}

func SetDefaultConfig(config *Config) {
	if len(config.Ui.Editor) <= 0 {
		config.Ui.Editor = "more"
	}
	if len(config.Core.Linefeed) <= 0 {
		config.Core.Linefeed = "unix"
	}
}

func LoadConfig() (Config, error) {
	var config Config
	var err error
	home, err := GetHomePath()
	if err == nil {
		configFile := filepath.Join(home, ".gookmarkrc")
		_, err = os.Stat(configFile)
		if err == nil {
			_, err = toml.DecodeFile(configFile, &config)
		}
		SetDefaultConfig(&config)
		err = nil
	}
	return config, err
}

func ShowCurrentConfig(config *Config) {
	fmt.Println("ui.editor=" + config.Ui.Editor)
	fmt.Println("core.linefeed=" + config.Core.Linefeed)
}

func SetNewConfig(config *Config, newConfig string) {
	buff := strings.SplitN(newConfig, ".", 2)
	if len(buff) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid option")
		return
	}
	section := buff[0]
	buff = strings.SplitN(buff[1], "=", 2)
	if len(buff) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid option")
		return
	}
	option := buff[0]
	value := buff[1]

	if section == "ui" {
		if option == "editor" {
			config.Ui.Editor = value
		}
	} else if section == "core" {
		if option == "linefeed" {
			config.Core.Linefeed = value
		}
	}
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	err := encoder.Encode(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	home, err := GetHomePath()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	configFile := filepath.Join(home, ".gookmarkrc")
	fp, err := os.OpenFile(configFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer fp.Close()
	fmt.Fprintln(fp, buffer.String())
}
