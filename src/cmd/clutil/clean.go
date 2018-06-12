package main

import "os"

func clean() error {
	d, err := getLanguageDirectory()

	if err != nil {
		return err
	}

	return os.RemoveAll(d)
}
