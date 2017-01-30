package rbt

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestSetInsertRemoveRandomly(t *testing.T) {
	s := NewSet(less)

	for i := 0; i < MAX_ITERS; i++ {
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
	s := NewSet(less)

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
