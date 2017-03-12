package core

type functionType struct {
	signature Signature
	function  func(...*Thunk) Value
}

func (f functionType) call(args Arguments) Value {
	ts, err := f.signature.Bind(args)

	if err != nil {
		return err
	}

	return f.function(ts...)
}

func NewLazyFunction(s Signature, f func(...*Thunk) Value) *Thunk {
	return Normal(functionType{
		signature: s,
		function:  f,
	})
}

func NewStrictFunction(s Signature, f func(...Value) Value) *Thunk {
	return NewLazyFunction(s, func(ts ...*Thunk) Value {
		for _, t := range ts {
			go t.Eval()
		}

		vs := make([]Value, len(ts))

		for i, t := range ts {
			vs[i] = t.Eval()

			if err, ok := vs[i].(ErrorType); ok {
				return err
			}
		}

		return f(vs...)
	})
}

func (f functionType) string() Value {
	return StringType("<function>")
}
