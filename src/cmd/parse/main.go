package main

import (
	"github.com/raviqqe/tisp/src/lib/parse"
	"github.com/docopt/docopt-go"
	"github.com/kr/pretty"
	"io/ioutil"
)

func main() {
	args := getArgs()

	source, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		panic(err.Error())
	}

	pretty.Println(parse.Parse(string(source)))
}

func getArgs() map[string]interface{} {
	usage := `Tisp parser

Usage:
  parse <filename>`

	args, err := docopt.Parse(usage, nil, true, "Tisp 0.0.0", false)

	if err != nil {
		panic(err.Error())
	}

	return args
}
