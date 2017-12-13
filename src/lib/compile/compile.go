package compile

import (
	"io/ioutil"
	"os"

	"github.com/coel-lang/coel/src/lib/desugar"
	"github.com/coel-lang/coel/src/lib/parse"
)

// Compile compiles a main module of a path into effects of thunks.
func Compile(path string) ([]Effect, error) {
	p, s, err := readFileOrStdin(path)

	if err != nil {
		return nil, err
	}

	m, err := parse.MainModule(p, s)

	if err != nil {
		return nil, err
	}

	c := newCompiler(builtinsEnvironment(), newModulesCache())
	return c.compileModule(desugar.Desugar(m))
}

func readFileOrStdin(path string) (string, string, error) {
	f := os.Stdin

	if path == "" {
		path = "<stdin>"
	} else {
		var err error
		f, err = os.Open(path)

		if err != nil {
			return "", "", err
		}
	}

	bs, err := ioutil.ReadAll(f)

	if err != nil {
		return "", "", err
	}

	return path, string(bs), nil
}
