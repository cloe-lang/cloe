package core

type collection interface {
	callable

	include(Value) Value
	index(Value) Value
	insert(Value, Value) Value
	merge(...Value) Value
	delete(Value) Value
	toList() Value
	size() Value
}

// Include returns true if a collection includes an element, or false otherwise.
var Include = NewStrictFunction(
	NewSignature([]string{"collection", "elem"}, "", nil, ""),
	func(vs ...Value) Value {
		c, err := evalCollection(vs[0])

		if err != nil {
			return err
		}

		return c.include(EvalPure(vs[1]))
	})

// Index extracts an element corresponding with a key.
var Index = NewStrictFunction(
	NewSignature([]string{"collection", "key"}, "keys", nil, ""),
	func(vs ...Value) Value {
		v := vs[0]
		l := cons(vs[1], vs[2])

		for !l.Empty() {
			c, err := evalCollection(v)

			if err != nil {
				return err
			}

			v = c.index(EvalPure(l.First()))

			if l, err = EvalList(l.Rest()); err != nil {
				return err
			}
		}

		return v
	})

// Insert inserts an element into a collection.
var Insert FunctionType

func initInsert() FunctionType {
	return NewLazyFunction(
		NewSignature([]string{"collection"}, "keyValuePairs", nil, ""),
		func(vs ...Value) (result Value) {
			c, err := evalCollection(vs[0])

			if err != nil {
				return err
			}

			l, err := EvalList(vs[1])

			if err != nil {
				return err
			}

			for !l.Empty() {
				k := l.First()
				l, err = EvalList(l.Rest())

				if err != nil {
					return err
				}

				c, err = evalCollection(c.insert(EvalPure(k), l.First()))

				if err != nil {
					return err
				}

				if l, err = EvalList(l.Rest()); err != nil {
					return err
				}
			}

			return c
		})
}

// Merge merges more than 2 collections.
var Merge FunctionType

func initMerge() FunctionType {
	return NewLazyFunction(
		NewSignature([]string{"collection"}, "collections", nil, ""),
		func(vs ...Value) Value {
			c, err := evalCollection(vs[0])

			if err != nil {
				return err
			}

			l, err := EvalList(vs[1])

			if err != nil {
				return err
			} else if l.Empty() {
				return c
			}

			vs = []Value{}

			for !l.Empty() {
				vs = append(vs, l.First())

				if l, err = EvalList(l.Rest()); err != nil {
					return err
				}
			}

			return c.merge(vs...)
		})
}

// Delete deletes an element corresponding with a key.
var Delete = NewStrictFunction(
	NewSignature([]string{"collection", "elem"}, "", nil, ""),
	func(vs ...Value) Value {
		c, err := evalCollection(vs[0])

		if err != nil {
			return err
		}

		return c.delete(EvalPure(vs[1]))
	})

// Size returns a size of a collection.
var Size = newUnaryCollectionFunction(func(c collection) Value { return c.size() })

// ToList converts a collection into a list of its elements.
var ToList = newUnaryCollectionFunction(func(c collection) Value { return c.toList() })

func newUnaryCollectionFunction(f func(c collection) Value) Value {
	return NewLazyFunction(
		NewSignature([]string{"collection"}, "", nil, ""),
		func(vs ...Value) Value {
			c, err := evalCollection(vs[0])

			if err != nil {
				return err
			}

			return f(c)
		})
}
