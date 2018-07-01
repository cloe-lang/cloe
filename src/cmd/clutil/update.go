package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloe-lang/cloe/src/lib/consts"
)

func update() error {
	d, err := consts.GetLanguageDirectory()

	if err != nil {
		return err
	}

	w := sync.WaitGroup{}

	rs, err := listRepositories(filepath.Join(d, "src"))

	if err != nil {
		return err
	}

	for _, r := range rs {
		w.Add(1)
		go func(r string) {
			defer w.Done()

			gitPull(filepath.Join(d, "src", r))
		}(r)
	}

	w.Wait()

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
