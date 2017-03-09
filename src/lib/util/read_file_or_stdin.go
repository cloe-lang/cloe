package util

import (
	"io/ioutil"
	"os"
)

// ReadFileOrStdin reads a file or stdin.
// maybeFilename should be string or nil.
func ReadFileOrStdin(maybeFilename interface{}) (string, string) {
	filename := "<stdin>"
	file := os.Stdin

	if s, ok := maybeFilename.(string); ok {
		filename = s

		var err error
		file, err = os.Open(s)

		if err != nil {
			Fail(err.Error())
		}
	}

	source, err := ioutil.ReadAll(file)

	if err != nil {
		Fail(err.Error())
	}

	return filename, string(source)
}
