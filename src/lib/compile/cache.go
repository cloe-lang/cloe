package compile

import (
	"errors"
	"path"

	"github.com/cloe-lang/cloe/src/lib/modules"
)

type modulesCache map[string]module

func newModulesCache() modulesCache {
	c := make(modulesCache, len(modules.Modules))

	for s, m := range modules.Modules {
		c[s] = m
	}

	return c
}

func (c modulesCache) Set(p string, m module) error {
	if !path.IsAbs(p) {
		return errors.New("module path is not absolute")
	}

	c[p] = m

	return nil
}

func (c modulesCache) Get(p string) (module, bool) {
	m, ok := c[p]
	return m, ok
}
