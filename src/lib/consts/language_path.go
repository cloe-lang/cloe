package consts

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cloe-lang/cloe/src/lib/utils"
)

func getLanguageSubDirectory(s string) (string, error) {
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

	d = filepath.Join(d, s)

	if err := utils.MkdirRecursively(d); err != nil {
		return "", err
	}

	return d, nil
}

// GetModulesDirectory gets a modules directory of the current user.
func GetModulesDirectory() (string, error) {
	return getLanguageSubDirectory("src")
}

// GetCommandsDirectory gets a commands directory of the current user.
func GetCommandsDirectory() (string, error) {
	return getLanguageSubDirectory("bin")
}
