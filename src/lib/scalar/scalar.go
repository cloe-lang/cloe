package scalar

import (
	"fmt"
	"strconv"

	"github.com/tisp-lang/tisp/src/lib/core"
)

var predefined = map[string]*core.Thunk{
	"true":       core.True,
	"false":      core.False,
	"nil":        core.Nil,
	"$emptyList": core.EmptyList,
	"$emptyDict": core.EmptyDictionary,
}

// Convert converts a string into a scalar value of a number, string, bool, or
// nil.
func Convert(name string) (*core.Thunk, error) {
	if t, ok := predefined[name]; ok {
		return t, nil
	}

	if n, err := strconv.ParseInt(name, 0, 64); err == nil && name[0] == '0' {
		return core.NewNumber(float64(n)), nil
	}

	if n, err := strconv.ParseFloat(name, 64); err == nil {
		return core.NewNumber(n), nil
	}

	if s, err := strconv.Unquote(name); err == nil {
		return core.NewString(s), nil
	}

	return nil, fmt.Errorf("the name, %s not found", name)
}

// Defined checks if a given name is a defined scalar value or not.
func Defined(name string) bool {
	if _, err := Convert(name); err == nil {
		return true
	}

	return false
}
