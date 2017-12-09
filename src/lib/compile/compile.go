package compile

import (
	"io/ioutil"
	"os"

	"github.com/tisp-lang/tisp/src/lib/desugar"
	"github.com/tisp-lang/tisp/src/lib/parse"
)

// Compile compiles a main module of a path into effects of thunks.
func Compile(path string) []Effect {
	module, err := parse.MainModule(readFileOrStdin(path))

	if err != nil {
		panic(err)
	}

	c := newCompiler(builtinsEnvironment(), newModulesCache())
	return c.compile(desugar.Desugar(module))
}

func (c *compiler) subModule(path string) module {
	module, err := parse.SubModule(readFileOrStdin(path))

	if err != nil {
		panic(err)
	}

	cc := newCompiler(builtinsEnvironment(), c.cache)
	c = &cc
	c.compile(desugar.Desugar(module))

	return c.env.toMap()
}

func readFileOrStdin(filename string) (string, string) {
	file := os.Stdin

	if filename == "" {
		filename = "<stdin>"
	} else {
		var err error
		file, err = os.Open(filename)

		if err != nil {
			panic(err)
		}
	}

	source, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	return filename, string(source)
}
