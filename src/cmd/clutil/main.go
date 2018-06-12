// +build go1.8

package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

func main() {
	args := getArgs()

	i, err := newInstaller(args["<repo>"].(string))

	if err != nil {
		fail(err)
	}

	err = i.InstallModule()

	if err != nil {
		fail(err)
	}

	err = i.InstallCommands()

	if err != nil {
		fail(err)
	}
}

func getArgs() map[string]interface{} {
	usage := `Cloe utility

Usage:
  clutil install <repo>

Options:
  -h, --help  Show this help.`

	args, err := docopt.ParseArgs(usage, os.Args[1:], "0.1.0")

	if err != nil {
		fail(err)
	} else if args["<repo>"] == nil {
		args["<repo>"] = ""
	}

	return args
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
