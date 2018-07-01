package main

import (
	"os"

	"github.com/cloe-lang/cloe/src/lib/consts"
)

func clean() error {
	for _, f := range [](func() (string, error)){consts.GetModulesDirectory, consts.GetCommandsDirectory} {
		d, err := f()

		if err != nil {
			return err
		}

		if err := os.RemoveAll(d); err != nil {
			return err
		}
	}

	return nil
}
