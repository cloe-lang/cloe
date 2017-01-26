package rbt

import "fmt"

type color bool

const (
	red   color = true
	black color = false
)

type Ordered interface {
	Less(Ordered) bool
}

type node struct {
	color       color
	value       Ordered
	left, right *node
}

func newNode(c color, o Ordered, l, r *node) *node {
	return &node{
		color: c,
		value: o,
		left:  l,
		right: r,
	}
}

func (n *node) insert(o Ordered) *node {
	m := *n.insertRed(o)
	m.color = black
	return &m
}

func (n *node) insertRed(o Ordered) *node {
	if n == nil {
		return newNode(red, o, nil, nil)
	}

	m := *n

	if n.value.Less(o) {
		m.left = m.left.insertRed(o)
	} else if o.Less(n.value) {
		m.right = m.right.insertRed(o)
	} else {
		return n
	}

	return m.balance()
}

func (n *node) balance() *node {
	if n.color == red {
		return n
	}

	newN := func(
		o Ordered,
		lo Ordered, ll, lr *node,
		ro Ordered, rl, rr *node) *node {
		return newNode(red, o, newNode(black, lo, ll, lr), newNode(black, ro, rl, rr))
	}

	l := n.left
	r := n.right

	if l != nil {
		lb := l != nil && l.color == red
		ll := l.left
		lr := l.right

		newLN := func(o, lo Ordered, ll, lr, rl *node) *node {
			return newN(o, lo, ll, lr, n.value, rl, r)
		}

		if lb && ll != nil && ll.color == red {
			return newLN(l.value, ll.value, ll.left, ll.right, lr)
		} else if lb && lr != nil && lr.color == red {
			return newLN(lr.value, l.value, ll, lr.left, lr.right)
		}
	} else if r != nil {
		rb := r != nil && r.color == red
		rl := r.left
		rr := r.right

		newRN := func(o, ro Ordered, lr, rl, rr *node) *node {
			return newN(o, n.value, l, lr, ro, rl, rr)
		}

		if rb && rl != nil && rl.color == red {
			return newRN(r.value, rr.value, rl, rr.left, rr.right)
		} else if rb && rr != nil && rr.color == red {
			return newRN(rl.value, r.value, rl.left, rl.right, rr)
		}
	}

	return n
}

func (n *node) search(o Ordered) (Ordered, bool) {
	if n == nil {
		return nil, false
	} else if n.value.Less(o) {
		return n.left.search(o)
	} else if o.Less(n.value) {
		return n.right.search(o)
	}

	return n.value, true
}

func (n *node) dump() {
	n.dumpWithIndent(0)
}

func (n *node) dumpWithIndent(i int) {
	for j := 0; j < i; j++ {
		fmt.Printf(" ")
	}

	if n == nil {
		fmt.Println(nil)
		return
	}

	fmt.Println(n.color, n.value)

	k := i + 2
	n.left.dumpWithIndent(k)
	n.right.dumpWithIndent(k)
}
