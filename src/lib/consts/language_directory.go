package consts

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cloe-lang/cloe/src/lib/utils"
)

// GetLanguageDirectory gets the language directory of the current user from
// proper environment variables.
func GetLanguageDirectory() (string, error) {
	d := os.Getenv("CLOE_PATH")

	if d == "" {
		h := os.Getenv("HOME")

		if h == "" {
			return "", errors.New("HOME environment variable is not set")
		}

		d = filepath.Join(h, ".cloe")
	}

	if !path.IsAbs(d) {
		return "", fmt.Errorf("language path, %s is not absolute", d)
	}

	if err := utils.MkdirRecursively(d); err != nil {
		return "", err
	}

	return d, nil
}
