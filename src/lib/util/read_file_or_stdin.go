package util

import (
	"io/ioutil"
	"os"
)

// ReadFileOrStdin reads a file or stdin.
// filename can be an empty string ("") to read stdin.
func ReadFileOrStdin(filename string) (string, string) {
	file := os.Stdin

	if filename == "" {
		filename = "<stdin>"
	} else {
		var err error
		file, err = os.Open(filename)

		if err != nil {
			PanicError(err)
		}
	}

	source, err := ioutil.ReadAll(file)

	if err != nil {
		PanicError(err)
	}

	return filename, string(source)
}
