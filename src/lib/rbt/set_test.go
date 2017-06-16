package rbt

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetInsertRemoveRandomly(t *testing.T) {
	s := NewSet(compare)

	for i := 0; i < MaxIters; i++ {
		x := generateKey()
		insert := rand.Int()%2 == 0

		if insert {
			s = s.Insert(x)
		} else {
			s = s.Remove(x)
		}

		ok := s.Include(x)
		assert.True(t, insert && ok || !insert && !ok)
	}
}

func TestSetFirstRest(t *testing.T) {
	xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := NewSet(compare)

	for _, x := range xs {
		s = s.Insert(x)
	}

	x, s := s.FirstRest()

	for _, expected := range xs {
		t.Log(x)
		assert.Equal(t, expected, x)
		x, s = s.FirstRest()
	}

	assert.Equal(t, nil, x)
	assert.True(t, s.Empty())
}

func TestSetMerge(t *testing.T) {
	NewSet(compare).Merge(NewSet(compare))
}
