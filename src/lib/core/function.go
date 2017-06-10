package core

import "github.com/tisp-lang/tisp/src/lib/systemt"

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

// NewLazyFunction creates a function whose arguments are evaluated lazily.
func NewLazyFunction(s Signature, f func(...*Thunk) Value) *Thunk {
	return Normal(functionType{
		signature: s,
		function:  f,
	})
}

// NewStrictFunction creates a function whose arguments are evaluated strictly.
func NewStrictFunction(s Signature, f func(...*Thunk) Value) *Thunk {
	return NewLazyFunction(s, func(ts ...*Thunk) Value {
		for _, t := range ts {
			tt := t
			systemt.Daemonize(func() { tt.Eval() })
		}

		return f(ts...)
	})
}

func (f functionType) string() Value {
	return StringType("<function>")
}
