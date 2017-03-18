package load

import (
	"github.com/raviqqe/tisp/src/lib/compile"
	"github.com/raviqqe/tisp/src/lib/desugar"
	"github.com/raviqqe/tisp/src/lib/parse"
	"github.com/raviqqe/tisp/src/lib/util"
)

// MainModule loads a main module as outputs.
func MainModule(filename string) []compile.Output {
	module, err := parse.Parse(util.ReadFileOrStdin(filename))

	if err != nil {
		util.Fail(err.Error())
	}

	return compile.Compile(desugar.Desugar(module))
}
