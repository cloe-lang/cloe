package compile

import "path/filepath"

type modulesCache map[string]module

func newModulesCache() modulesCache {
	return modulesCache{}
}

func (c modulesCache) Set(p string, m module) error {
	p, err := normalizePath(p)

	if err != nil {
		return err
	}

	c[p] = m

	return nil
}

func (c modulesCache) Get(p string) (module, bool, error) {
	p, err := normalizePath(p)

	if err != nil {
		return nil, false, err
	}

	m, ok := c[p]
	return m, ok, nil
}

func normalizePath(p string) (string, error) {
	return filepath.Abs(p)
}
