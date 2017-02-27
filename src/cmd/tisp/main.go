package main

import (
	"github.com/docopt/docopt-go"
	"github.com/raviqqe/tisp/src/lib/compile"
	"github.com/raviqqe/tisp/src/lib/parse"
	"github.com/raviqqe/tisp/src/lib/run"
	"io/ioutil"
	"log"
)

func main() {
	args := getArgs()

	source, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		log.Fatalln(err.Error())
	}

	module, err := parse.Parse(string(source))

	if err != nil {
		log.Fatalln(err.Error())
	}

	run.Run(compile.Compile(module))
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
