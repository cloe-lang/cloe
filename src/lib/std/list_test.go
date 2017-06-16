package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestList(t *testing.T) {
	n1, n2 := core.NewNumber(123), core.NewNumber(456)
	l := core.PApp(List, n1, n2)

	assert.Equal(t, n1.Eval(), core.PApp(core.First, l).Eval())
	assert.Equal(t, n2.Eval(), core.PApp(core.First, core.PApp(core.Rest, l)).Eval())
}
