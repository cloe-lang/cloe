package main

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cloe-lang/cloe/src/lib/consts"
)

type installer struct {
	url              *url.URL
	modulesDirectory string
	moduleDirectory  string
}

func newInstaller(s string) (installer, error) {
	d, err := consts.GetModulesDirectory()

	if err != nil {
		return installer{}, err
	}

	u, err := url.Parse(s)

	if err != nil {
		return installer{}, err
	}

	return installer{
		u,
		d,
		filepath.Join(d, u.Hostname(), filepath.FromSlash(u.Path)),
	}, nil
}

func (i installer) InstallModule() error {
	j, err := os.Stat(i.moduleDirectory)

	if err == nil {
		if !j.IsDir() {
			return errors.New("module directory is not a directory")
		}

		return gitPull(i.moduleDirectory)
	}

	if !os.IsNotExist(err) {
		return err
	}

	return gitClone(i.url.String(), i.moduleDirectory)
}

func (i installer) InstallCommands() error {
	ps := []string{}

	err := filepath.Walk(i.moduleDirectory, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !i.IsDir() && filepath.Base(p) == "main.cloe" {
			ps = append(ps, p)
		}

		return nil
	})

	if err != nil {
		return err
	}

	d, err := consts.GetCommandsDirectory()

	if err != nil {
		return err
	}

	for _, p := range ps {
		bs, err := ioutil.ReadFile(p)

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(d, filepath.Base(filepath.Dir(p))), bs, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}
