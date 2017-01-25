package rbt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type key int

func (k key) Less(o Ordered) bool {
	return k < o.(key)
}

func (n *node) dump(t *testing.T) {
	n.dumpWithIndent(t, 0)
}

func (n *node) dumpWithIndent(t *testing.T, i int) {
	for j := 0; j < i; j++ {
		fmt.Printf(" ")
	}

	if n == nil {
		fmt.Println(nil)
		return
	}

	fmt.Println(n.value)

	k := i + 2
	n.left.dumpWithIndent(t, k)
	n.right.dumpWithIndent(t, k)
}

func TestNode(t *testing.T) {
	k := key(3)

	n := (*node)(nil)
	t.Log(n)
	n = n.insert(k)
	t.Log(n)

	assert.Equal(t, n.search(k), k)
}

func TestNodeBalance(t *testing.T) {
	ks := []key{1, 2, 3, 4, 5, 6}
	n := (*node)(nil)

	for _, k := range ks {
		n.dump(t)
		n = n.insert(k)
	}

	n.dump(t)
}
