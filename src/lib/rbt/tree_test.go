package rbt

import (
	"math/rand"
	"testing"
)

func TestNodeInsertRemoveRandomly(t *testing.T) {
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
