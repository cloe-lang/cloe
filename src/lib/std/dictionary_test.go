package std

import (
	"testing"

	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestDictionary(t *testing.T) {
	for _, ts := range [][]*core.Thunk{
		{},
		{core.NewNumber(123), core.NewNumber(456)},
		{core.NewNumber(123), core.True, core.NewNumber(456), core.False},
	} {
		ks := make([]core.Object, len(ts)/2)
		vs := make([]*core.Thunk, len(ts)/2)

		for i, t := range ts {
			if i%2 == 0 {
				ks[i/2] = t.Eval()
			} else {
				vs[i/2] = t
			}
		}

		assert.True(t, bool(core.PApp(core.Equal,
			core.PApp(Dictionary, ts...),
			core.NewDictionary(ks, vs)).Eval().(core.BoolType)))
	}
}
