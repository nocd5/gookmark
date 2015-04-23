package main

import (
	"github.com/codegangsta/cli"
)

var Flags = []cli.Flag{
	cli.StringFlag{
		Name:  "group, g",
		Value: "bookmark",
		Usage: "bookmark group name",
	},
}
