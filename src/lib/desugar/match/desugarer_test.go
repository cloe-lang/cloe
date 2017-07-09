package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualPatterns(t *testing.T) {
	for _, p := range []interface{}{
		"42",
		"foo",
		app("$list"),
		app("$list", "true"),
		app("$list", app("$list")),
		app("$dict"),
		app("$dict", "123", "true"),
		app("$dict", app("$dict")),
	} {
		assert.True(t, equalPatterns(p, p))
	}

	for _, ps := range [][2]interface{}{
		{"42", "2049"},
		{"foo", "bar"},
		{app("$list"), app("$dict")},
		{app("$list", "true"), app("$list", "false")},
		{app("$list", app("$list")), app("$list", app("$list", "42"))},
		{app("$dict"), app("$dict", "0", "1")},
		{app("$dict", "123", "true"), app("$dict", "456", "true")},
		{app("$dict", app("$dict")), app("$dict", app("$list"))},
	} {
		assert.True(t, !equalPatterns(ps[0], ps[1]))
	}
}
