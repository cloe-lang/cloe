package compile

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestNewEffect(t *testing.T) {
	t.Log(NewEffect(core.Nil, false))
}

func TestEffectValue(t *testing.T) {
	assert.NotEqual(t, (*core.Thunk)(nil), NewEffect(core.Nil, false).Value())
}

func TestEffectExpanded(t *testing.T) {
	assert.False(t, NewEffect(core.Nil, false).Expanded())
	assert.True(t, NewEffect(core.EmptyList, true).Expanded())
}
