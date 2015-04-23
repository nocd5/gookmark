package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

func Contains(s []string, e string) bool {
	for _, _s := range s {
		if _s == e {
			return true
		}
	}
	return false
}

func ReplaceArgsByAlias(src []string) []string {
	var dest []string
	var subcommands []string

	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return src
	}

	for _, cmd := range Commands {
		subcommands = append(subcommands, cmd.Name)
	}

	skipNext := false
	for _, arg := range src {
		if skipNext {
			dest = append(dest, arg)
			skipNext = false
			continue
		}

		if Contains(subcommands, arg) {
			skipNext = true
		}

		alias := config.Alias[arg]
		if len(alias) > 0 {
			dest = append(dest, strings.SplitN(alias, " ", 2)...)
		} else {
			dest = append(dest, arg)
		}
	}

	return dest
}

func main() {
	app := cli.NewApp()
	app.Name = "gookmark"
	app.Version = Version
	app.Usage = ""
	app.Author = "nocd5"
	app.Email = "nocd5rd@gmail.com"
	app.Commands = Commands

	app.Run(ReplaceArgsByAlias(os.Args))
}
