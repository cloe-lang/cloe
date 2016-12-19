package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `Risp interpreter

Usage:
  risp <filename>

Options:
  -h --help     Show this help.`

	args, _ := docopt.Parse(usage, nil, true, "Risp 0.0.0", false)
	fmt.Println(args)
}
