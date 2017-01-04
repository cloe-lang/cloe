package main

import (
	_ "../../lib/parse"
	"../../lib/vm"
	"fmt"
	"github.com/docopt/docopt-go"
)

func main() {
	args := getArgs()

	fmt.Println(args)

	vm.Main().Eval()
}

func getArgs() map[string]interface{} {
	usage := `Risp interpreter

Usage:
  risp <filename>

Options:
  -h --help     Show this help.`

	args, _ := docopt.Parse(usage, nil, true, "Risp 0.0.0", false)

	return args
}
