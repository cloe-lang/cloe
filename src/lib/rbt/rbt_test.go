package rbt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type key int

func (k key) Less(o Ordered) bool {
	return k < o.(key)
}

func TestNode(t *testing.T) {
	k := key(3)

	n := (*node)(nil)
	t.Log(n)
	n = n.insert(k)
	t.Log(n)

	assert.Equal(t, n.search(k), k)
}

func TestTreeConstruction(t *testing.T) {
	ks := []key{1, 2, 3}
	n := (*node)(nil)

	for _, k := range ks {
		t.Log(n)
		n = n.insert(k)
	}
}
