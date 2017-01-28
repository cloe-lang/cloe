package rbt

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestTreeInsertRemoveRandomly(t *testing.T) {
	tr := NewTree(less)

	for i := 0; i < MAX_ITERS; i++ {
		x := generateKey()
		insert := rand.Int()%2 == 0

		if insert {
			tr = tr.Insert(x)
		} else {
			tr, _ = tr.Remove(x)
		}

		_, ok := tr.Search(x)

		if insert && !ok || !insert && ok {
			t.Fail()
		}
	}
}

func TestTreeFirstRest(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	tr := NewTree(less)

	for _, x := range xs {
		tr = tr.Insert(x)
	}

	x, f := tr.FirstRest()

	for _, xpected := range xs {
		t.Log(x)
		assert.Equal(t, xpected, x)
		x, f = f()
	}

	assert.Equal(t, nil, x)
	assert.Equal(t, FirstRestFunc(nil), f)
}
