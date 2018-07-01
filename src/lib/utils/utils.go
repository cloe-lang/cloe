package utils

import (
	"fmt"
	"os"
)

// MkdirRecursively works just like mkdir -p.
func MkdirRecursively(d string) error {
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
