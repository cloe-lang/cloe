package main

import (
	"fmt"
	"os"

	"github.com/cloe-lang/cloe/src/lib/consts"
)

func mkdirp(d string) error {
	if err := os.Mkdir(d, 0700); err != nil && !os.IsExist(err) {
		return err
	}

	if i, err := os.Stat(d); err != nil {
		return err
	} else if !i.IsDir() {
		return fmt.Errorf("%s is not a directory", d)
	}

	return nil
}

func getLanguageDirectory() (string, error) {
	d := os.Getenv(consts.PathName)

	if d == "" {
		return "", fmt.Errorf("%v environment variable is not set", consts.PathName)
	}

	err := mkdirp(d)

	if err != nil {
		return "", err
	}

	return d, nil
}
