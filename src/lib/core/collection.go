package core

type collection interface {
	callable

	include(Value) Value
	index(Value) Value
	insert(*Thunk, *Thunk) Value
	merge(...*Thunk) Value
	delete(Value) Value
	toList() Value
	size() Value
}

// Include returns true if a collection includes an element, or false otherwise.
var Include = NewStrictFunction(
	NewSignature([]string{"collection", "elem"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		c, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		v = ts[1].Eval()
		if err, ok := v.(ErrorType); ok {
			return err
		}

		return c.include(v)
	})

// Index extracts an element corresponding with a key.
var Index = NewStrictFunction(
	NewSignature([]string{"collection", "key"}, nil, "keys", nil, nil, ""),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l := cons(ts[1], ts[2])

		for !l.Empty() {
			c, ok := v.(collection)

			if !ok {
				return NotCollectionError(v)
			}

			v = ensureNormal(c.index(l.First().Eval()))

			var err Value
			l, err = l.Rest().EvalList()

			if err != nil {
				return err
			}
		}

		return v
	})

// Insert inserts an element into a collection.
var Insert = NewLazyFunction(
	NewSignature([]string{"collection"}, nil, "keyValuePairs", nil, nil, ""),
	func(ts ...*Thunk) (result Value) {
		v := ts[0].Eval()
		c, ok := v.(collection)

		if !ok {
			return NotCollectionError(v)
		}

		l := ts[1]

		for {
			if v := ReturnIfEmptyList(l, c); v != nil {
				return v
			}

			x := PApp(First, l)
			l = PApp(Rest, l)
			y := PApp(First, l)
			l = PApp(Rest, l)

			v = ensureNormal(c.insert(x, y))
			c, ok = v.(collection)

			if !ok {
				return NotCollectionError(v)
			}
		}
	})

// Merge merges more than 2 collections.
var Merge = NewLazyFunction(
	NewSignature([]string{"collection"}, nil, "collections", nil, nil, ""),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		c, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		l, err := ts[1].EvalList()

		if err != nil {
			return err
		} else if l.Empty() {
			return c
		}

		ts = []*Thunk{}

		for !l.Empty() {
			ts = append(ts, l.First())

			l, err = l.Rest().EvalList()

			if err != nil {
				return err
			}
		}

		return c.merge(ts...)
	})

// Delete deletes an element corresponding with a key.
var Delete = NewStrictFunction(
	NewSignature([]string{"collection", "elem"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		d, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return d.delete(ts[1].Eval())
	})

// Size returns a size of a collection.
var Size = newUnaryCollectionFunction(func(c collection) Value { return c.size() })

// ToList converts a collection into a list of its elements.
var ToList = newUnaryCollectionFunction(func(c collection) Value { return c.toList() })

func newUnaryCollectionFunction(f func(c collection) Value) *Thunk {
	return NewLazyFunction(
		NewSignature([]string{"collection"}, nil, "", nil, nil, ""),
		func(ts ...*Thunk) Value {
			v := ts[0].Eval()
			c, ok := v.(collection)

			if !ok {
				return NotCollectionError(v)
			}

			return f(c)
		})
}
