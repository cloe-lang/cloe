package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloe-lang/cloe/src/lib/consts"
	git "gopkg.in/src-d/go-git.v4"
)

func update() error {
	d, err := consts.GetModulesDirectory()

	if err != nil {
		return err
	}

	rs, err := listRepositories(d)

	if err != nil {
		return err
	}

	for _, r := range rs {
		if err := gitPull(r); err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	}

	return nil
}

func listRepositories(root string) ([]string, error) {
	rs := []string{}

	ds, err := ioutil.ReadDir(root)

	if err != nil {
		return nil, err
	}

	for _, d := range ds {
		d := filepath.Join(root, d.Name())

		if i, err := os.Stat(filepath.Join(d, ".git")); err == nil && i.IsDir() {
			rs = append(rs, d)
			continue
		}

		ss, err := listRepositories(d)

		if err != nil {
			return nil, err
		}

		rs = append(rs, ss...)
	}

	return rs, nil
}
