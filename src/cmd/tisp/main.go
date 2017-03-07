package main

import (
	"github.com/docopt/docopt-go"
	"github.com/raviqqe/tisp/src/lib/compile"
	"github.com/raviqqe/tisp/src/lib/desugar"
	"github.com/raviqqe/tisp/src/lib/parse"
	"github.com/raviqqe/tisp/src/lib/run"
	"log"
)

func main() {
	args := getArgs()

	module, err := parse.Parse(args["<filename>"].(string))

	if err != nil {
		log.Fatalln(err.Error())
	}

	run.Run(compile.Compile(desugar.Desugar(module)))
}

func getArgs() map[string]interface{} {
	usage := `Tisp interpreter

Usage:
  tisp <filename>

Options:
  -h --help     Show this help.`

	args, _ := docopt.Parse(usage, nil, true, "Tisp 0.0.0", false)

	return args
}
