package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gookmark"
	app.Version = Version
	app.Usage = ""
	app.Author = "nocd5"
	app.Email = "nocd5rd@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
