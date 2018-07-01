package main

import (
	"os"

	"github.com/cloe-lang/cloe/src/lib/consts"
)

func clean() error {
	d, err := consts.GetLanguageDirectory()

	if err != nil {
		return err
	}

	return os.RemoveAll(d)
}
