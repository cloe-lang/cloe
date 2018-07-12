// +build go1.8

package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
)

func main() {
	if err := command(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}

func command() error {
	c, args, err := getArgs()

	if err != nil {
		return err
	}

	switch c {
	case "install":
		return install(args["<repo>"].(string))
	case "update":
		return update()
	case "clean":
		return clean()
	default:
		panic("invalid subcommand")
	}
}

func getArgs() (string, map[string]interface{}, error) {
	usage := `Cloe utility

Usage:
	clutil install <repo>
	clutil update
	clutil clean

Options:
	-h, --help  Show this help.`

	args, err := docopt.ParseArgs(usage, os.Args[1:], "0.1.0")

	if err != nil {
		return "", nil, err
	} else if args["<repo>"] == nil {
		args["<repo>"] = ""
	}

	return os.Args[1], args, nil
}
