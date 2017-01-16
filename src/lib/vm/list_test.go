package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListEqual(t *testing.T) {
	for _, tss := range [][][]*Thunk{
		{{}, {}},
		{{True}, {True}},
		{{True, False}, {True, False}},
	} {
		assert.True(t, testEqual(NewList(tss[0]...), NewList(tss[1]...)))
	}

	for _, tss := range [][][]*Thunk{
		{{}, {True}},
		{{True}, {False}},
		{{True, True}, {True, True, True}},
	} {
		assert.True(t, !testEqual(NewList(tss[0]...), NewList(tss[1]...)))
	}
}
