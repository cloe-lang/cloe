package compile

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/cloe-lang/cloe/src/lib/desugar"
	"github.com/cloe-lang/cloe/src/lib/parse"
)

// Compile compiles a main module of a path into effects of thunks.
func Compile(p string) ([]Effect, error) {
	q, s, err := readFileOrStdin(p)

	if err != nil {
		return nil, err
	}

	m, err := parse.MainModule(q, s)

	if err != nil {
		return nil, err
	}

	c := newCompiler(builtinsEnvironment(), newModulesCache())
	return c.compileModule(desugar.Desugar(m), filepath.ToSlash(path.Dir(p))) // path.Dir("") == "."
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
