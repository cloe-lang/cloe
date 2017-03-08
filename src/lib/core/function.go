package core

type functionType struct {
	signature Signature
	function  func(...*Thunk) Object
}

func (f functionType) call(args Arguments) Object {
	ts, err := f.signature.Bind(args)

	if err != nil {
		return err
	}

	return f.function(ts...)
}

func NewLazyFunction(s Signature, f func(...*Thunk) Object) *Thunk {
	return Normal(functionType{
		signature: s,
		function:  f,
	})
}

func NewStrictFunction(s Signature, f func(...Object) Object) *Thunk {
	return NewLazyFunction(s, func(ts ...*Thunk) Object {
		for _, t := range ts {
			go t.Eval()
		}

		os := make([]Object, len(ts))

		for i, t := range ts {
			os[i] = t.Eval()

			if err, ok := os[i].(ErrorType); ok {
				return err
			}
		}

		return f(os...)
	})
}

func (f functionType) string() Object {
	return StringType("<function>")
}
