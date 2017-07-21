package compile

import (
	"io/ioutil"
	"os"

	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/desugar"
	"github.com/tisp-lang/tisp/src/lib/parse"
)

// MainModule compiles a main module of a path into outputs of thunks.
func MainModule(path string) []Output {
	module, err := parse.MainModule(readFileOrStdin(path))

	if err != nil {
		panic(err)
	}

	c := newCompiler()
	return c.compile(desugar.Desugar(module))
}

// SubModule compiles a sub module of a path into a map of names to thunks.
func SubModule(path string) map[string]*core.Thunk {
	return subModule(readFileOrStdin(path))
}

func subModule(filename, source string) map[string]*core.Thunk {
	module, err := parse.SubModule(filename, source)

	if err != nil {
		panic(err)
	}

	c := newCompiler()
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
