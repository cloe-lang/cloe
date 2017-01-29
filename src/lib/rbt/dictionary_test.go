package rbt

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestDictionaryInsertRemoveRandomly(t *testing.T) {
	d := NewDictionary(less)

	for i := 0; i < MAX_ITERS; i++ {
		k := generateKey()
		insert := rand.Int()%2 == 0

		if insert {
			d = d.Insert(k, keyToValue(k))
		} else {
			d = d.Remove(k)
		}

		_, ok := d.Search(k)

		if insert && !ok || !insert && ok {
			t.Fail()
		}
	}
}

func TestDictionaryFirstRest(t *testing.T) {
	ks := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	d := NewDictionary(less)

	for _, k := range ks {
		d = d.Insert(k, keyToValue(k))
	}

	k, v, d := d.FirstRest()

	for _, expected := range ks {
		t.Log(k, v)
		assert.Equal(t, expected, k)
		assert.Equal(t, keyToValue(expected), v)
		k, v, d = d.FirstRest()
	}

	assert.Equal(t, nil, k)
	assert.Equal(t, nil, v)
	assert.True(t, d.Empty())
}

func keyToValue(k int) int {
	return k + 1000
}
