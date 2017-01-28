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

func TestNodeInsertRandomly(t *testing.T) {
	n := (*node)(nil)

	for i := 0; i < MAX_ITERS; i++ {
		k := key(rand.Int() % MAX_KEY)
		old := n
		n = n.insert(k)

		n.rank() // check ranks

		if !n.checkColors() {
			fmt.Println("KEY:", k)
			fmt.Println("OLD:")
			old.dump()
			fmt.Println("NEW:")
			n.dump()
			t.FailNow()
		}
	}
}

func TestNodeRemoveRandomly(t *testing.T) {
	n := (*node)(nil)

	for i := 0; i < MAX_ITERS; i++ {
		k := key(rand.Int() % MAX_KEY)
		old := n

		remove := rand.Int()%3 == 0

		if remove {
			n, _ = n.remove(k)

			_, ok := n.search(k)

			if ok {
				t.Fail()
			}
		} else {
			n = n.insert(k)

			_, ok := n.search(k)

			if !ok {
				t.Fail()
			}
		}

		n.rank() // check ranks

		if !n.checkColors() {
			if remove {
				fmt.Println("REMOVE")
			} else {
				fmt.Println("INSERT")
			}

			fmt.Println("KEY:", k)
			fmt.Println("OLD:")
			old.dump()
			fmt.Println("NEW:")
			n.dump()
			t.FailNow()
		}
	}
}
