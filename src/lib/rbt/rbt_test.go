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

	kk, ok := n.search(k)
	assert.True(t, ok)
	assert.Equal(t, kk, k)
}

func TestNodeBalance(t *testing.T) {
	ks := []key{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n.dump()
		n = n.insert(k)
	}

	n.dump()
}

func TestNodeTakeMax(t *testing.T) {
	ks := []key{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n = n.insert(k)
		n.dump()

		o, m, _ := n.takeMax()
		assert.Equal(t, k, o.(key))
		m.dump()
	}
}

func TestNodeRemove(t *testing.T) {
	ks := []key{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n = n.insert(k)
	}

	n.dump()

	for _, k := range ks {
		n, _ = n.remove(k)
		n.dump()
	}

	n.dump()
}
