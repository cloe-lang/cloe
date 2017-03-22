package core

type collection interface {
	callable

	include(Value) Value
	index(Value) Value
	insert(...*Thunk) Value
	merge(...*Thunk) Value
	delete(Value) Value
	toList() Value
	size() Value
}

// Include returns true if a collection includes an element, or false otherwise.
var Include = NewStrictFunction(
	NewSignature(
		[]string{"collection", "elem"}, nil, "",
		nil, nil, "",
	),
	func(vs ...Value) Value {
		c, ok := vs[0].(collection)

		if !ok {
			return TypeError(vs[0], "collection")
		}

		if err, ok := vs[1].(ErrorType); ok {
			return err
		}

		return c.include(vs[1])
	})

// Index extracts an element corresponding with a key.
var Index = NewStrictFunction(
	NewSignature(
		[]string{"collection", "key"}, nil, "",
		nil, nil, "",
	),
	func(vs ...Value) Value {
		i, ok := vs[0].(collection)

		if !ok {
			return TypeError(vs[0], "collection")
		}

		return i.index(vs[1])
	})

// Insert inserts an element into a collection.
var Insert = NewLazyFunction(
	NewSignature(
		[]string{"collection"}, nil, "values",
		nil, nil, "",
	),
	func(ts ...*Thunk) (result Value) {
		v := ts[0].Eval()
		c, ok := v.(collection)

		if !ok {
			return NotCollectionError(v)
		}

		v = ts[1].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToThunks()

		if err != nil {
			return err
		}

		return c.insert(ts...)
	})

// Merge merges 2 collections.
var Merge = NewLazyFunction(
	NewSignature(
		[]string{"x"}, nil, "ys",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		m, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		v = ts[1].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToThunks()

		if err != nil {
			return err
		}

		if len(ts) == 0 {
			return m
		}

		return m.merge(ts...)
	})

// Delete deletes an element corresponding with a key.
var Delete = NewStrictFunction(
	NewSignature(
		[]string{"collection", "elem"}, nil, "",
		nil, nil, "",
	),
	func(vs ...Value) Value {
		d, ok := vs[0].(collection)

		if !ok {
			return TypeError(vs[0], "collection")
		}

		return d.delete(vs[1])
	})

// Size returns a size of a collection.
var Size = NewStrictFunction(
	NewSignature(
		[]string{"collection"}, nil, "",
		nil, nil, "",
	),
	func(vs ...Value) Value {
		v := vs[0]
		c, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return c.size()
	})

// ToList converts a collection into a list of its elements.
var ToList = NewStrictFunction(
	NewSignature(
		[]string{"listLike"}, nil, "",
		nil, nil, "",
	),
	func(vs ...Value) Value {
		v := vs[0]
		l, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return l.toList()
	})
