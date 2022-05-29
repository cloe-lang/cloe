package scalar

import (
	"fmt"
	"strconv"

	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/cloe-lang/cloe/src/lib/core"
)

var predefined = func() map[string]core.Value {
	m := map[string]core.Value{
		"true":                       core.True,
		"false":                      core.False,
		"nil":                        core.Nil,
		consts.Names.EmptyList:       core.EmptyList,
		consts.Names.EmptyDictionary: core.EmptyDictionary,
	}

	for k, v := range m {
		if k[:1] != "$" {
			m["$"+k] = v
		}
	}

	return m
}()

// Convert converts a string into a scalar value of a number, string, bool, or
// nil.
func Convert(name string) (core.Value, error) {
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
