package core

import "github.com/coel-lang/coel/src/lib/systemt"

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
	return Normal(functionType{s, f})
}

// NewStrictFunction creates a function whose arguments are evaluated strictly.
func NewStrictFunction(s Signature, f func(...*Thunk) Value) *Thunk {
	return NewLazyFunction(s, func(ts ...*Thunk) Value {
		systemt.Daemonize(func() {
			for _, t := range ts {
				tt := t
				systemt.Daemonize(func() { tt.Eval() })
			}
		})

		return f(ts...)
	})
}

// NewEffectFunction creates a effect function which returns an effect value.
func NewEffectFunction(s Signature, f func(...*Thunk) Value) *Thunk {
	return Normal(functionType{s, func(ts ...*Thunk) Value { return newEffect(Normal(f(ts...))) }})
}

func (f functionType) string() Value {
	return StringType("<function>")
}
