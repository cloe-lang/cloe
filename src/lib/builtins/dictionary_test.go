package builtins

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestDictionary(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{},
		{core.NewNumber(123), core.NewNumber(456)},
		{core.NewNumber(123), core.True, core.NewNumber(456), core.False},
	} {
		kvs := make([]core.KeyValue, 0, len(ts)/2)

		for i := 0; i < len(ts); i += 2 {
			kvs = append(kvs, core.KeyValue{ts[i], ts[i+1]})
		}

		assert.True(t, bool(core.PApp(core.Equal,
			core.PApp(Dictionary, ts...),
			core.NewDictionary(kvs)).Eval().(core.BoolType)))
	}
}
