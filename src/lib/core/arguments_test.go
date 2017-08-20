package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArgumentsMerging(t *testing.T) {
	NewArguments([]PositionalArgument{
		NewPositionalArgument(Nil, false),
		NewPositionalArgument(EmptyList, true),
		NewPositionalArgument(Nil, false),
		NewPositionalArgument(EmptyList, true),
	}, nil, nil)
}

func TestArgumentsEmpty(t *testing.T) {
	for _, a := range []Arguments{
		NewArguments(nil, nil, []*Thunk{Nil}),
		NewArguments(nil, nil, []*Thunk{NewDictionary([]Value{Nil.Eval()}, []*Thunk{Nil})}),
	} {
		th := a.empty()
		assert.NotEqual(t, (*Thunk)(nil), th)
		v := th.Eval()
		t.Log(v)
		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}
