package match

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/stretchr/testify/assert"
)

func TestEqualPatterns(t *testing.T) {
	dictionary := consts.Names.DictionaryFunction
	list := consts.Names.ListFunction

	for _, p := range []interface{}{
		"42",
		"foo",
		app(list),
		app(list, "true"),
		app(list, app(list)),
		app(dictionary),
		app(dictionary, "123", "true"),
		app(dictionary, app(dictionary)),
	} {
		assert.True(t, equalPatterns(p, p))
	}

	for _, ps := range [][2]interface{}{
		{"42", "2049"},
		{"foo", "bar"},
		{app(list), app(dictionary)},
		{app(list, "true"), app(list, "false")},
		{app(list, app(list)), app(list, app(list, "42"))},
		{app(dictionary), app(dictionary, "0", "1")},
		{app(dictionary, "123", "true"), app(dictionary, "456", "true")},
		{app(dictionary, app(dictionary)), app(dictionary, app(list))},
	} {
		assert.True(t, !equalPatterns(ps[0], ps[1]))
	}
}
