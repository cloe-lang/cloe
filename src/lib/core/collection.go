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
	NewSignature(
		[]string{"collection", "key"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		i, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return i.index(ts[1].Eval())
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
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		d, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return d.delete(ts[1].Eval())
	})

// Size returns a size of a collection.
var Size = NewLazyFunction(
	NewSignature(
		[]string{"collection"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		c, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return c.size()
	})

// ToList converts a collection into a list of its elements.
var ToList = NewLazyFunction(
	NewSignature(
		[]string{"listLike"}, nil, "",
		nil, nil, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return l.toList()
	})
