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

	if l != nil && l.color == red {
		ll := l.left
		lr := l.right

		newLN := func(o, lo Ordered, ll, lr, rl *node) *node {
			return newN(o, lo, ll, lr, n.value, rl, r)
		}

		if ll != nil && ll.color == red {
			return newLN(l.value, ll.value, ll.left, ll.right, lr)
		} else if lr != nil && lr.color == red {
			return newLN(lr.value, l.value, ll, lr.left, lr.right)
		}
	} else if r != nil && r.color == red {
		rl := r.left
		rr := r.right

		newRN := func(o, ro Ordered, lr, rl, rr *node) *node {
			return newN(o, n.value, l, lr, ro, rl, rr)
		}

		if rl != nil && rl.color == red {
			return newRN(r.value, rr.value, rl, rr.left, rr.right)
		} else if rr != nil && rr.color == red {
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

func (n *node) remove(o Ordered) (*node, bool) {
	_, ok := n.search(o)

	if !ok {
		return n, false
	}

	n, _ = n.removeOne(o)
	m := *n
	m.color = black
	return &m, true
}

func (n *node) removeOne(o Ordered) (*node, bool) {
	if n == nil {
		return nil, true
	} else if n.value.Less(o) {
		n, balanced := n.left.removeOne(o)

		if balanced {
			return n, true
		}

		return n.balanceLeft()
	} else if o.Less(n.value) {
		n, balanced := n.right.removeOne(o)

		if balanced {
			return n, true
		}

		return n.balanceRight()
	}

	if n.left == nil {
		return n.right, n.color == red
	}

	o, l, balanced := n.takeMax()

	m := newNode(n.color, o, l, n.right)

	if balanced {
		return m, true
	}

	return m.balanceLeft()
}

func (n *node) takeMax() (Ordered, *node, bool) {
	if n.right == nil {
		return n.value, n.left, n.color == red
	}

	o, r, balanced := n.right.takeMax()

	m := *n
	m.right = r

	if balanced {
		return o, &m, true
	}

	n, balanced = m.balanceRight()
	return o, n, balanced
}

func (n *node) balanceLeft() (*node, bool) {
	if n.right.color == red && n.color == red {
		panic("What the hell!")
	}

	if n.right.color == red {
		l, _ := newNode(red, n.value, n.left, n.right.left).balanceLeft()
		return newNode(black, n.right.value, l, n.right.right), true
	}

	if n.right.left.color == red {
		return newNode(
			n.color,
			n.right.left.value,
			newNode(black, n.value, n.left, n.right.left.left),
			newNode(black, n.right.value, n.right.left.right, n.right.right)), true
	} else if n.right.right.color == red {
		r := *n.right.right
		r.color = black

		return newNode(
			n.color,
			n.right.value,
			newNode(black, n.value, n.left, n.right.left),
			&r), true
	}

	r := *n.right
	r.color = red

	m := *n
	m.color = black
	m.right = &r

	return &m, n.color == black
}

func (n *node) balanceRight() (*node, bool) {
	if n.right.color == red && n.color == red {
		panic("What the hell!")
	}

	if n.right.color == red {
		l, _ := newNode(red, n.value, n.left, n.right.left).balanceLeft()
		return newNode(black, n.right.value, l, n.right.right), true
	}

	if n.right.left.color == red {
		return newNode(
			n.color,
			n.right.left.value,
			newNode(black, n.value, n.left, n.right.left.left),
			newNode(black, n.right.value, n.right.left.right, n.right.right)), true
	} else if n.right.right.color == red {
		r := *n.right.right
		r.color = black

		return newNode(
			n.color,
			n.right.value,
			newNode(black, n.value, n.left, n.right.left),
			&r), true
	}

	r := *n.right
	r.color = red

	m := *n
	m.color = black
	m.left = &r

	return &m, n.color == black
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
