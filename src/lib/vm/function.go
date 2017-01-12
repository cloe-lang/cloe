package vm

type functionType func(...*Thunk) Object

func (f functionType) Call(ts ...*Thunk) Object {
	return f(ts...)
}

func NewLazyFunction(f func(...*Thunk) Object) *Thunk {
	return Normal(functionType(f))
}

func NewStrictFunction(f func(...Object) Object) *Thunk {
	return NewLazyFunction(func(ts ...*Thunk) Object {
		for _, t := range ts {
			go t.Eval()
		}

		os := make([]Object, len(ts))

		for i, t := range ts {
			os[i] = t.Eval()
		}

		return f(os...)
	})
}
