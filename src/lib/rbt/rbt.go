package rbt

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

func (n *node) black() *node {
	m := *n
	m.color = black
	return &m
}

func (n *node) insert(o Ordered) *node {
	if n == nil {
		return newNode(red, o, nil, nil) // black?
	}

	m := *n

	if n.value.Less(o) {
		m.left = m.left.insert(o)
	} else if o.Less(n.value) {
		m.right = m.right.insert(o)
	} else {
		return n
	}

	return m.balance()
}

func (n *node) balance() *node {
	if n.color == red {
		return n
	}

	l := n.left
	lb := l != nil && l.color == red
	ll := l.left
	lr := l.right

	r := n.right
	rb := r != nil && r.color == red
	rl := r.left
	rr := r.right

	newN := func(
		o Ordered,
		lo Ordered, ll, lr *node,
		ro Ordered, rl, rr *node) *node {
		return newNode(red, o, newNode(black, lo, ll, lr), newNode(black, ro, rl, rr))
	}

	newRN := func(o, lo Ordered, ll, lr, rl *node) *node {
		return newN(o, lo, ll, lr, n.value, rl, r)
	}

	newLN := func(o, ro Ordered, lr, rl, rr *node) *node {
		return newN(o, n.value, l, lr, ro, rl, rr)
	}

	if lb && ll != nil && ll.color == red {
		return newRN(l.value, ll.value, ll.left, ll.right, lr)
	} else if lb && lr != nil && lr.color == red {
		return newRN(lr.value, l.value, ll, lr.left, lr.right)
	} else if rb && rl != nil && rl.color == red {
		return newLN(r.value, rr.value, rl, rr.left, rr.right)
	} else if rb && rr != nil && rr.color == red {
		return newLN(rl.value, r.value, rl.left, rl.right, rr)
	}

	return n
}
