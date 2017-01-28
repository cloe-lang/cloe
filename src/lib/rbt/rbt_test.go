package rbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
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
}

const (
	MAX_ITERS = 2000
	MAX_KEY   = MAX_ITERS / 2
)

func generateKey() key {
	return key(rand.Int() % MAX_KEY)
}

func TestNodeInsertRandomly(t *testing.T) {
	n := (*node)(nil)

	for i := 0; i < MAX_ITERS; i++ {
		k := generateKey()
		old := n

		n = n.insert(k)

		n.rank() // check ranks

		if !n.checkColors() {
			failWithDump(t, true, k, old, n)
		}
	}
}

func TestNodeRemoveRandomly(t *testing.T) {
	n := (*node)(nil)

	for i := 0; i < MAX_ITERS; i++ {
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

	for i := 0; i < MAX_ITERS; i++ {
		k := generateKey()
		old := n
		oldCopy := n.deepCopy()

		n, insert := n.insertOrRemove(t, k)

		n.rank() // check ranks

		if !(old.totalEqual(oldCopy) && n.checkColors()) {
			failWithDump(t, insert, k, old, n)
		}
	}
}

func (n *node) insertOrRemove(t *testing.T, x Ordered) (*node, bool) {
	insert := rand.Int()%2 == 0

	if insert {
		n = n.insert(x)
	} else {
		n, _ = n.remove(x)
	}

	_, ok := n.search(x)

	if insert && !ok || !insert && ok {
		t.Fail()
	}

	return n, insert
}

func failWithDump(t *testing.T, insert bool, k key, old, new *node) {
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
