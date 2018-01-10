package core

import "github.com/coel-lang/coel/src/lib/systemt"

// FunctionType represents a function.
type FunctionType func(Arguments) Value

func (f FunctionType) call(args Arguments) Value {
	return f(args)
}

func (f FunctionType) eval() Value {
	return f
}

// NewRawFunction creates a function which takes arguments directly.
func NewRawFunction(f func(Arguments) Value) FunctionType {
	return FunctionType(f)
}

// NewLazyFunction creates a function whose arguments are evaluated lazily.
func NewLazyFunction(s Signature, f func(...Value) Value) FunctionType {
	return NewRawFunction(func(args Arguments) Value {
		vs, err := s.Bind(args)

		if err != nil {
			return err
		}

		return f(vs...)
	})
}

// NewStrictFunction creates a function whose arguments are evaluated strictly.
func NewStrictFunction(s Signature, f func(...Value) Value) FunctionType {
	return NewLazyFunction(s, func(vs ...Value) Value {
		systemt.Daemonize(func() {
			for _, t := range vs {
				tt := t
				systemt.Daemonize(func() { tt.eval() })
			}
		})

		return f(vs...)
	})
}

// NewEffectFunction creates a effect function which returns an effect value.
func NewEffectFunction(s Signature, f func(...Value) Value) Value {
	ff := NewLazyFunction(s, f)

	return NewRawFunction(func(args Arguments) Value {
		return newEffect(App(ff, args))
	})
}

// Partial creates a partially-applied function with arguments.
var Partial = FunctionType(func(vars Arguments) Value {
	return NewRawFunction(func(args Arguments) Value {
		vars := vars
		v := EvalPure(vars.nextPositional())
		f, ok := v.(callable)

		if !ok {
			return NotCallableError(v)
		}

		return f.call(vars.Merge(args))
	})
})

func (f FunctionType) string() Value {
	return NewString("<function>")
}
