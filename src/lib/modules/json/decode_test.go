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
		t.Log(core.PApp(decode, core.NewString(s)).Eval())
	}
}

func TestDecodeWithNonString(t *testing.T) {
	_, ok := core.PApp(decode, core.Nil).Eval().(core.ErrorType)
	assert.True(t, ok)
}

func TestDecodeString(t *testing.T) {
	for _, s := range jsons {
		t.Log(core.PApp(core.ToString, decodeString(s)).Eval())
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
		_, ok := decodeString(s).Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}

func TestConvertToValueWithInvalidType(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	convertToValue(0)
}
