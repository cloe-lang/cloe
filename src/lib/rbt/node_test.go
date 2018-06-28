package rbt

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func compare(x, y interface{}) int {
	return x.(int) - y.(int)
}

func TestNode(t *testing.T) {
	k := 3

	n := (*node)(nil)
	t.Log(n)
	n = n.insert(k, compare)
	t.Log(n)

	kk, ok := n.search(k, compare)
	assert.True(t, ok)
	assert.Equal(t, kk, k)
}

func TestNodeBalance(t *testing.T) {
	ks := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n.dump()
		n = n.insert(k, compare)
	}

	n.dump()
}

func TestNodeTakeMax(t *testing.T) {
	ks := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k1 := range ks {
		n = n.insert(k1, compare)
		n.dump()

		k2, m, _ := n.takeMax()
		assert.Equal(t, k1, k2.(int))
		m.dump()
	}
}

func TestNodeRemove(t *testing.T) {
	ks := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n = n.insert(k, compare)
	}

	n.dump()

	for _, k := range ks {
		n = n.remove(k, compare)
		n.dump()
	}
}

const (
	MaxIters = 256
	MaxKey   = MaxIters / 2
)

func generateKey() int {
	return rand.Int() % MaxKey
}

func TestNodeInsertRandomly(t *testing.T) {
	n := (*node)(nil)

	for i := 0; i < MaxIters; i++ {
		k := generateKey()
		old := n

		n = n.insert(k, compare)

		n.rank() // check ranks

		if !n.checkColors() {
			failWithDump(t, true, k, old, n)
		}
	}
}

func TestNodeRemoveRandomly(t *testing.T) {
	n := (*node)(nil)

	for i := 0; i < MaxIters; i++ {
		k := generateKey()
		old := n

		var insert bool
		n, insert = n.insertOrRemove(t, k)

		n.dump()
		n.rank() // check ranks

		if !n.checkColors() {
			failWithDump(t, insert, k, old, n)
		}
	}
}

func TestInsertRemovePersistency(t *testing.T) {
	n := (*node)(nil)

	for i := 0; i < MaxIters; i++ {
		k := generateKey()
		old := n
		oldCopy := n.deepCopy()

		var insert bool
		n, insert = n.insertOrRemove(t, k)

		n.dump()
		n.rank() // check ranks

		if !(old.equal(oldCopy) && n.checkColors()) {
			failWithDump(t, insert, k, old, n)
		}
	}
}

func (n *node) insertOrRemove(t *testing.T, x interface{}) (*node, bool) {
	insert := rand.Int()%2 == 0

	if insert {
		n = n.insert(x, compare)
	} else {
		n = n.remove(x, compare)
	}

	_, ok := n.search(x, compare)
	assert.True(t, insert && ok || !insert && !ok)

	return n, insert
}

func failWithDump(t *testing.T, insert bool, k int, old, new *node) {
	if insert {
		fmt.Println("INSERT")
	} else {
		fmt.Println("REMOVE")
	}

	fmt.Println("KEY:", k)
	fmt.Println("OLD:")
	old.dump()
	fmt.Println("NEW:")
	new.dump()
	t.FailNow()
}

func TestNodeEqual(t *testing.T) {
	n0 := (*node)(nil)
	n1 := n0.insert(1, compare)
	n2 := n0.insert(2, compare)
	n3 := n1.insert(2, compare)
	n4 := n3.insert(3, compare)

	for _, test := range []struct {
		n1, n2 *node
		equal  bool
	}{
		{n0, n0, true},
		{n1, n1, true},
		{n2, n2, true},
		{n3, n3, true},
		{n4, n4, true},
		{n0, n1, false},
		{n0, n2, false},
		{n0, n3, false},
		{n0, n4, false},
		{n1, n2, false},
		{n1, n3, false},
		{n2, n3, false},
		{n2, n4, false},
		{n3, n4, false},
	} {
		assert.Equal(t, test.equal, test.n1.equal(test.n2))
	}
}

func TestNodeRankPanic(t *testing.T) {
	assert.Panics(t, func() {
		n := (*node)(nil).insert(0, compare).insert(1, compare).insert(2, compare)
		n.dump()
		n.left.color = red
		n.dump()

		n.rank()
	})
}

func TestNodeCheckColorsError(t *testing.T) {
	n := (*node)(nil).insert(0, compare).insert(1, compare).insert(2, compare).insert(3, compare)
	n.dump()
	n.right.color = red
	n.dump()
	assert.True(t, !n.checkColors())
}
