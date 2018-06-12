package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/pkg/errors"
)

type installer struct {
	url             *url.URL
	languagePath    string
	moduleDirectory string
}

func newInstaller(s string) (installer, error) {
	p := os.Getenv(consts.PathName)

	if p == "" {
		return installer{}, errors.Errorf("%v environment variable is not set", consts.PathName)
	}

	u, err := url.Parse(s)

	if err != nil {
		return installer{}, err
	}

	return installer{
		u,
		p,
		filepath.Join(p, "src", u.Hostname(), filepath.FromSlash(u.Path)),
	}, nil
}

func (i installer) InstallModule() error {
	j, err := os.Stat(i.moduleDirectory)

	if err == nil {
		if !j.IsDir() {
			return errors.New("module directory is not a directory")
		}

		if err := os.Chdir(i.moduleDirectory); err != nil {
			return err
		}

		return exec.Command("git", "pull").Run()
	}

	if !os.IsNotExist(err) {
		return err
	}

	return exec.Command("git", "clone", i.url.String(), i.moduleDirectory).Run()
}

func (i installer) InstallCommands() error {
	ps := []string{}

	err := filepath.Walk(i.moduleDirectory, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if i, err := os.Stat(p); err == nil && !i.IsDir() && filepath.Base(p) == "main.cloe" {
			ps = append(ps, p)
		}

		return nil
	})

	if err != nil {
		return err
	}

	d := filepath.Join(i.languagePath, "bin")

	if err = os.Mkdir(d, 0700); err != nil && !os.IsExist(err) {
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
