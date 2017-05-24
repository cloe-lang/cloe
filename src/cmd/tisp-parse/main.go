package main

import (
	"github.com/docopt/docopt-go"
	"github.com/kr/pretty"
	"github.com/tisp-lang/tisp/src/lib/parse"
	"github.com/tisp-lang/tisp/src/lib/util"
)

func main() {
	pretty.Println(parse.MainModule(util.ReadFileOrStdin(getArgs()["<filename>"].(string))))
}

func getArgs() map[string]interface{} {
	usage := `Tisp parser

Usage:
  parse [<filename>]`

	args, err := docopt.Parse(usage, nil, true, "Tisp 0.0.0", false)

	if err != nil {
		util.Fail(err.Error())
	} else if args["<filename>"] == nil {
		args["<filename>"] = ""
	}

	return args
}
