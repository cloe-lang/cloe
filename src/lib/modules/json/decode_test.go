package json

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

var jsons = []string{
	`true`,
	`false`,
	`123`,
	`-3.14`,
	`"foo"`,
	`null`,
	`[]`,
	`["foo", 42, "bar"]`,
	`{}`,
	`{"foo": 42, "bar": true, "baz": "blah"}`,
}

func TestDecode(t *testing.T) {
	for _, s := range jsons {
		t.Log(core.EvalPure(core.PApp(decode, core.NewString(s))))
	}
}

func TestDecodeWithNonString(t *testing.T) {
	_, ok := core.EvalPure(core.PApp(decode, core.Nil)).(core.ErrorType)
	assert.True(t, ok)
}

func TestDecodeString(t *testing.T) {
	for _, s := range jsons {
		t.Log(core.EvalPure(core.PApp(core.ToString, decodeString(s))))
	}
}

func TestDecodeStringWithInvalidJSON(t *testing.T) {
	for _, s := range []string{
		`s`,
		`[`,
		`{`,
		`nul`,
		`nil`,
	} {
		_, ok := core.EvalPure(decodeString(s)).(core.ErrorType)
		assert.True(t, ok)
	}
}

func TestConvertToValueWithInvalidType(t *testing.T) {
	assert.Panics(t, func() {
		convertToValue(0)
	})
}
