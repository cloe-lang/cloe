// +build go1.8

package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/cloe-lang/cloe/src/lib/compile"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/cloe-lang/cloe/src/lib/run"
	"github.com/docopt/docopt-go"
)

func main() {
	args := getArgs()

	if args["--debug"].(bool) {
		debug.Debug = true
	} else {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case error:
					printToStderr(x.Error())
				case string:
					printToStderr(x)
				default:
					panic(x)
				}

				os.Exit(1)
			}
		}()
	}

	if args["--profile"] != nil {
		f, err := os.Create(args["--profile"].(string))
		if err != nil {
			panic(err)
		}

		err = pprof.StartCPUProfile(f)
		if err != nil {
			panic(err)
		}

		defer pprof.StopCPUProfile()
	}

	es, err := compile.Compile(args["<filename>"].(string))

	if err != nil {
		panic(err)
	}

	run.Run(es)
}

func getArgs() map[string]interface{} {
	usage := `Cloe interpreter

Usage:
  cloe [-d] [-p <filename>] [<filename>]

Options:
  -d, --debug  Turn on debug mode.
  -p, --profile <filename>  Turn on profiling.
  -h, --help  Show this help.`

	args, err := docopt.ParseArgs(usage, os.Args[1:], "0.1.0")

	if err != nil {
		printToStderr(err.Error())
		os.Exit(1)
	} else if args["<filename>"] == nil {
		args["<filename>"] = ""
	}

	return args
}

func printToStderr(s string) {
	if _, err := fmt.Fprintf(os.Stderr, strings.TrimSpace(s)+"\n"); err != nil {
		panic(err)
	}
}
