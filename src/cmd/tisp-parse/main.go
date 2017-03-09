package main

import (
	"github.com/docopt/docopt-go"
	"github.com/kr/pretty"
	"github.com/raviqqe/tisp/src/lib/parse"
	"github.com/raviqqe/tisp/src/lib/util"
)

func main() {
	pretty.Println(parse.Parse(util.ReadFileOrStdin(getArgs()["<filename>"])))
}

func getArgs() map[string]interface{} {
	usage := `Tisp parser

Usage:
  parse [<filename>]`

	args, err := docopt.Parse(usage, nil, true, "Tisp 0.0.0", false)

	if err != nil {
		panic(err.Error())
	}

	return args
}
