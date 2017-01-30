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
		old := tr
		insert := rand.Int()%2 == 0

		if insert {
			tr = tr.Insert(x)
		} else {
			tr = tr.Remove(x)
		}

		_, ok := old.Search(x)

		if insert && !ok {
			assert.Equal(t, 1, tr.Size()-old.Size())
		} else if insert && ok {
			assert.Equal(t, tr.Size(), old.Size())
		} else if !insert && ok {
			assert.Equal(t, 1, old.Size()-tr.Size())
		} else {
			assert.Equal(t, tr.Size(), old.Size())
		}

		_, ok = tr.Search(x)
		assert.True(t, insert && ok || !insert && !ok)
	}
}

func TestTreeFirstRest(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	tr := NewTree(less)

	for _, x := range xs {
		tr = tr.Insert(x)
	}

	x, tr := tr.FirstRest()

	for _, xpected := range xs {
		t.Log(x)
		assert.Equal(t, xpected, x)
		x, tr = tr.FirstRest()
	}

	assert.Equal(t, nil, x)
	assert.True(t, tr.Empty())
}

func TestTreeSize(t *testing.T) {
	tr := NewTree(less)
	assert.Equal(t, 0, tr.Size())
	tr = tr.Insert(0)
	assert.Equal(t, 1, tr.Size())
	tr = tr.Insert(1)
	assert.Equal(t, 2, tr.Size())
	tr = tr.Insert(1)
	assert.Equal(t, 2, tr.Size())
	tr = tr.Remove(1)
	assert.Equal(t, 1, tr.Size())
	tr = tr.Remove(0)
	assert.Equal(t, 0, tr.Size())
	tr = tr.Remove(0)
	assert.Equal(t, 0, tr.Size())
}
