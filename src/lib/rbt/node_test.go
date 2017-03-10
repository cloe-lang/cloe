package rbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func less(x1, x2 interface{}) bool {
	i1, i2 := x1.(int), x2.(int)
	return i1 < i2
}

func TestNode(t *testing.T) {
	k := 3

	n := (*node)(nil)
	t.Log(n)
	n = n.insert(k, less)
	t.Log(n)

	kk, ok := n.search(k, less)
	assert.True(t, ok)
	assert.Equal(t, kk, k)
}

func TestNodeBalance(t *testing.T) {
	ks := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n.dump()
		n = n.insert(k, less)
	}

	n.dump()
}

func TestNodeTakeMax(t *testing.T) {
	ks := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n = n.insert(k, less)
		n.dump()

		o, m, _ := n.takeMax()
		assert.Equal(t, k, o.(int))
		m.dump()
	}
}

func TestNodeRemove(t *testing.T) {
	ks := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n := (*node)(nil)

	for _, k := range ks {
		n = n.insert(k, less)
	}

	n.dump()

	for _, k := range ks {
		n = n.remove(k, less)
		n.dump()
	}
}

const (
	MaxIters = 2000
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

		n = n.insert(k, less)

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

		n, insert := n.insertOrRemove(t, k)

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

		n, insert := n.insertOrRemove(t, k)

		n.rank() // check ranks

		if !(old.equal(oldCopy) && n.checkColors()) {
			failWithDump(t, insert, k, old, n)
		}
	}
}

func (n *node) insertOrRemove(t *testing.T, x interface{}) (*node, bool) {
	insert := rand.Int()%2 == 0

	if insert {
		n = n.insert(x, less)
	} else {
		n = n.remove(x, less)
	}

	_, ok := n.search(x, less)
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
