package main

import "os"

func clean() error {
	p, err := getLanguagePath()

	if err != nil {
		return err
	}

	return os.RemoveAll(p)
}
