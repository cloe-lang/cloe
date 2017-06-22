// +build go1.8

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/tisp-lang/tisp/src/lib/compile"
	"github.com/tisp-lang/tisp/src/lib/run"
)

func main() {
	args := getArgs()

	if !args["--debug"].(bool) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case error:
					printToStderr(x.Error())
				case string:
					printToStderr(x)
				default:
					if x != nil {
						panic(x)
					}
				}

				os.Exit(1)
			}
		}()
	}

	run.Run(compile.MainModule(args["<filename>"].(string)))
}

func getArgs() map[string]interface{} {
	usage := `Tisp interpreter

Usage:
  tisp [-d] [<filename>]

Options:
  -d, --debug  Turn on debug mode.
  -h, --help  Show this help.`

	args, err := docopt.Parse(usage, nil, true, "Tisp 0.0.0", false)

	if err != nil {
		printToStderr(err.Error())
		os.Exit(1)
	} else if args["<filename>"] == nil {
		args["<filename>"] = ""
	}

	return args
}

func printToStderr(s string) {
	fmt.Fprintf(os.Stderr, strings.TrimSpace(s)+"\n")
}
