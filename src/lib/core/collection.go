package core

// TODO: Create collection interface integrating some existing interfaces with
// methods of index, insert, merge, delete, size (or len?), include and toList.
// It should be implemented by StringType, ListType, DictionaryType and SetType.

type collection interface {
	index(Value) Value
	merge(ts ...*Thunk) Value
	delete(Value) Value
	toList() Value
}

// Index extracts an element corresponding with a key.
var Index = NewStrictFunction(
	NewSignature(
		[]string{"collection", "key"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		i, ok := vs[0].(collection)

		if !ok {
			return TypeError(vs[0], "collection")
		}

		return i.index(vs[1])
	})

// Merge merges 2 collections.
var Merge = NewLazyFunction(
	NewSignature(
		[]string{"x"}, []OptionalArgument{}, "ys",
		[]string{}, []OptionalArgument{}, "",
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
		[]string{"collection", "elem"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		d, ok := vs[0].(collection)

		if !ok {
			return TypeError(vs[0], "collection")
		}

		return d.delete(vs[1])
	})

// ToList converts a collection into a list of its elements.
var ToList = NewStrictFunction(
	NewSignature(
		[]string{"listLike"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		v := vs[0]
		l, ok := v.(collection)

		if !ok {
			return TypeError(v, "collection")
		}

		return l.toList()
	})
