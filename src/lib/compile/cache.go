package compile

import (
	"errors"
	"path"
)

type modulesCache map[string]module

func newModulesCache() modulesCache {
	return modulesCache{}
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
