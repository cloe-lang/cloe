// +build go1.8

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/tisp-lang/tisp/src/lib/compile"
	"github.com/tisp-lang/tisp/src/lib/run"
	"github.com/tisp-lang/tisp/src/lib/util"
)

func main() {
	defer func() {
		switch x := recover().(type) {
		case error:
			printToStderr(x.Error())
		case string:
			printToStderr(x)
		default:
			if x != nil {
				panic(x)
			}
		}
	}()

	run.Run(compile.MainModule(getArgs()["<filename>"].(string)))
}

func getArgs() map[string]interface{} {
	usage := `Tisp interpreter

Usage:
  tisp [<filename>]

Options:
  -h, --help  Show this help.`

	args, err := docopt.Parse(usage, nil, true, "Tisp 0.0.0", false)

	if err != nil {
		util.PanicError(err)
	} else if args["<filename>"] == nil {
		args["<filename>"] = ""
	}

	return args
}

func printToStderr(s string) {
	fmt.Fprintf(os.Stderr, strings.TrimSpace(s)+"\n")
}
