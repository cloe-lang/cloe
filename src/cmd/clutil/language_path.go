package main

import (
	"os"

	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/pkg/errors"
)

func getLanguagePath() (string, error) {
	p := os.Getenv(consts.PathName)

	if p == "" {
		return "", errors.Errorf("%v environment variable is not set", consts.PathName)
	}

	return p, nil
}
