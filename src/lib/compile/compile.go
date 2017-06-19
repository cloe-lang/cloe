package compile

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/desugar"
	"github.com/tisp-lang/tisp/src/lib/parse"
	"github.com/tisp-lang/tisp/src/lib/util"
)

// MainModule compiles a main module of a path into outputs of thunks.
func MainModule(path string) []Output {
	os := make([]Output, 0)
	p := parse.NewMainModuleParser(util.ReadFileOrStdin(path))
	c := newCompiler()

	for {
		if p.Finished() {
			return os
		}

		s, err := p.Parse()
		if err != nil {
			util.PanicError(err)
		}

		for _, s := range desugar.Desugar(s) {
			if o, ok := c.compile(s); ok {
				os = append(os, o)
			}
		}
	}
}

// SubModule compiles a sub module of a path into a map of names to thunks.
func SubModule(path string) map[string]*core.Thunk {
	p := parse.NewSubModuleParser(util.ReadFileOrStdin(path))
	c := newCompiler()

	for {
		if p.Finished() {
			return c.env.toMap()
		}

		s, err := p.Parse()
		if err != nil {
			util.PanicError(err)
		}

		for _, s := range desugar.Desugar(s) {
			c.compile(s)
		}
	}
}
