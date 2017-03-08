package main

import (
	"github.com/docopt/docopt-go"
	"github.com/kr/pretty"
	"github.com/raviqqe/tisp/src/lib/parse"
)

func main() {
	pretty.Println(parse.Parse(getArgs()["<filename>"].(string)))
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
