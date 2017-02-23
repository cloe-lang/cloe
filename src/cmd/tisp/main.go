package main

import (
	"../../lib/compile"
	"../../lib/parse"
	"../../lib/run"
	"github.com/docopt/docopt-go"
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
