package core

type sequence interface {
	collection

	insert(Value, Value) Value
}

// Insert inserts an element into a sequence.
var Insert FunctionType

func initInsert() FunctionType {
	return NewLazyFunction(
		NewSignature([]string{"collection"}, "indexValuePairs", nil, ""),
		func(vs ...Value) (result Value) {
			s, err := evalSequence(vs[0])

			if err != nil {
				return err
			}

			l, err := EvalList(vs[1])

			if err != nil {
				return err
			}

			for !l.Empty() {
				k := l.First()

				if l, err = EvalList(l.Rest()); err != nil {
					return err
				}

				s, err = evalSequence(s.insert(EvalPure(k), l.First()))

				if err != nil {
					return err
				}

				if l, err = EvalList(l.Rest()); err != nil {
					return err
				}
			}

			return s
		})
}
