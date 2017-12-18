package rbt

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictionaryInsertRemoveRandomly(t *testing.T) {
	d := NewDictionary(compare)

	for i := 0; i < MaxIters; i++ {
		k := generateKey()
		insert := rand.Int()%2 == 0

		if insert {
			d = d.Insert(k, keyToValue(k))
		} else {
			d = d.Remove(k)
		}

		_, ok := d.Search(k)
		assert.True(t, insert && ok || !insert && !ok)
	}
}

func TestDictionaryFirstRest(t *testing.T) {
	ks := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	d := NewDictionary(compare)

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

	assert.Nil(t, k)
	assert.Nil(t, v)
	assert.True(t, d.Empty())
}

func keyToValue(k int) int {
	return k + 1000
}

func TestDictionaryMerge(t *testing.T) {
	NewDictionary(compare).Merge(NewDictionary(compare))
}
