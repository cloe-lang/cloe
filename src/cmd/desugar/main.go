package main

import (
	"../../lib/desugar"
	"../../lib/parse"
	"fmt"
	"github.com/docopt/docopt-go"
	"io/ioutil"
)

func main() {
	args := getArgs()

	source, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		panic(err.Error())
	}

	fmt.Println(desugar.Desugar(parse.Parse(string(source))))
}

func getArgs() map[string]interface{} {
	usage := `Risp desugarer

Usage:
  desugar <filename>`

	args, err := docopt.Parse(usage, nil, true, "Risp 0.0.0", false)

	if err != nil {
		panic(err.Error())
	}

	return args
}
