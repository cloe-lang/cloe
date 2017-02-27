package std

import (
	. "github.com/raviqqe/tisp/src/lib/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList(t *testing.T) {
	n1, n2 := NewNumber(123), NewNumber(456)
	l := PApp(List, n1, n2)

	assert.Equal(t, n1.Eval(), PApp(First, l).Eval())
	assert.Equal(t, n2.Eval(), PApp(First, PApp(Rest, l)).Eval())
}
