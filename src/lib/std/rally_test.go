package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/systemt"
)

func TestRally(t *testing.T) {
	go systemt.RunDaemons()

	ts := []*core.Thunk{core.True, core.False, core.Nil}

	l1 := core.NewList(ts...)
	l2 := core.App(Rally, core.NewArguments(
		[]core.PositionalArgument{core.NewPositionalArgument(l1, true)},
		nil,
		nil))

	for i := 0; i < len(ts); i++ {
		e := core.PApp(core.First, l2)
		t.Logf("%#v", e.Eval())
		assert.True(t, bool(core.PApp(core.Include, l1, e).Eval().(core.BoolType)))

		var n *core.Thunk
		l1tmp := l1

		for i := 0; ; i++ {
			ee := core.PApp(core.First, l1tmp)
			t.Logf("%#v", ee.Eval())

			if core.PApp(core.Equal, e, ee).Eval().(core.BoolType) {
				n = core.NewNumber(float64(i))
				break
			} else if i >= len(ts) {
				t.FailNow()
			}

			l1tmp = core.PApp(core.Rest, l1tmp)
		}

		l1 = core.PApp(core.Delete, l1, n)
		l2 = core.PApp(core.Rest, l2)
	}

	assert.True(t, bool(core.PApp(core.Equal, core.EmptyList, l1).Eval().(core.BoolType)))
	assert.True(t, bool(core.PApp(core.Equal, core.EmptyList, l2).Eval().(core.BoolType)))
}
