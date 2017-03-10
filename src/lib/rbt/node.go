package rbt

import (
	"fmt"

	"github.com/raviqqe/tisp/src/lib/util"
)

type color bool

const (
	red   color = true
	black color = false
)

type node struct {
	color       color
	value       interface{}
	left, right *node
}

func newNode(c color, x interface{}, l, r *node) *node {
	return &node{
		color: c,
		value: x,
		left:  l,
		right: r,
	}
}

func (n *node) paint(c color) *node {
	m := *n
	m.color = c
	return &m
}

func (n *node) min() interface{} {
	if n == nil {
		return nil
	}

	if n.left == nil {
		return n.value
	}

	return n.left.min()
}

func (n *node) insert(x interface{}, less func(interface{}, interface{}) bool) *node {
	return n.insertRed(x, less).paint(black)
}

func (n *node) insertRed(x interface{}, less func(interface{}, interface{}) bool) *node {
	if n == nil {
		return newNode(red, x, nil, nil)
	}

	m := *n

	if less(x, n.value) {
		m.left = m.left.insertRed(x, less)
	} else if less(n.value, x) {
		m.right = m.right.insertRed(x, less)
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
		x interface{},
		lx interface{}, ll, lr *node,
		rx interface{}, rl, rr *node) *node {
		return newNode(red, x, newNode(black, lx, ll, lr), newNode(black, rx, rl, rr))
	}

	l := n.left
	r := n.right

	if l != nil && l.color == red {
		ll := l.left
		lr := l.right

		newLN := func(x, lx interface{}, ll, lr, rl *node) *node {
			return newN(x, lx, ll, lr, n.value, rl, r)
		}

		if ll != nil && ll.color == red {
			return newLN(l.value, ll.value, ll.left, ll.right, lr)
		} else if lr != nil && lr.color == red {
			return newLN(lr.value, l.value, ll, lr.left, lr.right)
		}
	}

	if r != nil && r.color == red {
		rl := r.left
		rr := r.right

		newRN := func(x, rx interface{}, lr, rl, rr *node) *node {
			return newN(x, n.value, l, lr, rx, rl, rr)
		}

		if rl != nil && rl.color == red {
			return newRN(rl.value, r.value, rl.left, rl.right, rr)
		} else if rr != nil && rr.color == red {
			return newRN(r.value, rr.value, rl, rr.left, rr.right)
		}
	}

	return n
}

func (n *node) search(x interface{}, less func(interface{}, interface{}) bool) (interface{}, bool) {
	if n == nil {
		return nil, false
	} else if less(x, n.value) {
		return n.left.search(x, less)
	} else if less(n.value, x) {
		return n.right.search(x, less)
	}

	return n.value, true
}

func (n *node) remove(x interface{}, less func(interface{}, interface{}) bool) *node {
	n, _ = n.removeOne(x, less)

	if n == nil {
		return nil
	}

	return n.paint(black)
}

func (n *node) removeOne(x interface{}, less func(interface{}, interface{}) bool) (*node, bool) {
	if n == nil {
		return nil, true
	} else if less(x, n.value) {
		l, balanced := n.left.removeOne(x, less)
		m := *n
		m.left = l

		if balanced {
			return &m, true
		}

		return m.balanceLeft()
	} else if less(n.value, x) {
		r, balanced := n.right.removeOne(x, less)
		m := *n
		m.right = r

		if balanced {
			return &m, true
		}

		return m.balanceRight()
	}

	if n.left == nil {
		return n.right, n.color == red
	}

	x, l, balanced := n.left.takeMax()

	m := newNode(n.color, x, l, n.right)

	if balanced {
		return m, true
	}

	return m.balanceLeft()
}

func (n *node) takeMax() (interface{}, *node, bool) {
	if n.right == nil {
		return n.value, n.left, n.color == red
	}

	x, r, balanced := n.right.takeMax()

	m := *n
	m.right = r

	if balanced {
		return x, &m, true
	}

	n, balanced = m.balanceRight()
	return x, n, balanced
}

func (n *node) balanceLeft() (*node, bool) {
	if n.right.color == red {
		l, _ := newNode(red, n.value, n.left, n.right.left).balanceLeft()
		return newNode(black, n.right.value, l, n.right.right), true
	}

	if n.right.left != nil && n.right.left.color == red {
		return newNode(
			n.color,
			n.right.left.value,
			newNode(black, n.value, n.left, n.right.left.left),
			newNode(black, n.right.value, n.right.left.right, n.right.right)), true
	} else if n.right.right != nil && n.right.right.color == red {
		return newNode(
			n.color,
			n.right.value,
			newNode(black, n.value, n.left, n.right.left),
			n.right.right.paint(black)), true
	}

	m := n.paint(black)
	m.right = n.right.paint(red)

	return m, n.color == red
}

func (n *node) balanceRight() (*node, bool) {
	if n.left.color == red {
		r, _ := newNode(red, n.value, n.left.right, n.right).balanceRight()
		return newNode(black, n.left.value, n.left.left, r), true
	}

	if n.left.right != nil && n.left.right.color == red {
		return newNode(
			n.color,
			n.left.right.value,
			newNode(black, n.left.value, n.left.left, n.left.right.left),
			newNode(black, n.value, n.left.right.right, n.right)), true
	} else if n.left.left != nil && n.left.left.color == red {
		return newNode(
			n.color,
			n.left.value,
			n.left.left.paint(black),
			newNode(black, n.value, n.left.right, n.right)), true
	}

	m := n.paint(black)
	m.left = n.left.paint(red)

	return m, n.color == red
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

	j := i + 2
	n.right.dumpWithIndent(j)
	n.left.dumpWithIndent(j)
}

func (n *node) checkColors() bool {
	if n == nil {
		return true
	}

	if n.color == red &&
		((n.left != nil && n.left.color == red) ||
			(n.right != nil && n.right.color == red)) {
		return false
	}

	if n.left.checkColors() && n.right.checkColors() {
		return true
	}

	return false
}

func (n *node) rank() int {
	if n == nil {
		return 1
	}

	if n.left.rank() != n.right.rank() {
		util.Fail("Red-black tree is unbalanced!")
	}

	r := 0

	if n.color == black {
		r = 1
	}

	return n.left.rank() + r
}

func (n *node) deepCopy() *node {
	if n == nil {
		return nil
	}

	return newNode(n.color, n.value, n.left.deepCopy(), n.right.deepCopy())
}

func (n *node) size() int {
	if n == nil {
		return 0
	}

	return n.left.size() + 1 + n.right.size()
}

func (n *node) equal(m *node) bool {
	if n == nil && m == nil {
		return true
	} else if n == nil || m == nil {
		return false
	}

	return n.color == m.color &&
		n.value == m.value &&
		n.left.equal(m.left) &&
		n.right.equal(m.right)
}
