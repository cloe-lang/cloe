package ast

import (
	"fmt"
)

// Convert converts an AST with a given function.
// This function visits all expressions and statements at least.
func Convert(f func(interface{}) interface{}, x interface{}) interface{} {
	if y := f(x); y != nil {
		return y
	}

	convert := func(x interface{}) interface{} {
		return Convert(f, x)
	}

	switch x := x.(type) {
	case string:
		return x
	case AnonymousFunction:
		return NewAnonymousFunction(convert(x.Signature()).(Signature), convert(x.Body()))
	case App:
		return NewApp(convert(x.Function()), convert(x.Arguments()).(Arguments), x.DebugInfo())
	case Arguments:
		ps := make([]PositionalArgument, 0, len(x.Positionals()))

		for _, p := range x.Positionals() {
			ps = append(ps, convert(p).(PositionalArgument))
		}

		ks := make([]KeywordArgument, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, convert(k).(KeywordArgument))
		}

		return NewArguments(ps, ks)
	case Import:
		return x
	case KeywordArgument:
		return NewKeywordArgument(x.Name(), convert(x.Value()))
	case DefFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			ls = append(ls, convert(l))
		}

		return NewDefFunction(
			x.Name(),
			convert(x.Signature()).(Signature),
			ls,
			convert(x.Body()),
			x.DebugInfo())
	case LetVar:
		return NewLetVar(x.Name(), convert(x.Expr()))
	case Match:
		cs := make([]MatchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			cs = append(cs, convert(c).(MatchCase))
		}

		return NewMatch(convert(x.Value()), cs)
	case MatchCase:
		return NewMatchCase(x.Pattern(), convert(x.Value()))
	case MutualRecursion:
		fs := make([]DefFunction, 0, len(x.DefFunctions()))

		for _, f := range x.DefFunctions() {
			fs = append(fs, convert(f).(DefFunction))
		}

		return NewMutualRecursion(fs, x.DebugInfo())
	case OptionalParameter:
		return NewOptionalParameter(x.Name(), convert(x.DefaultValue()))
	case Effect:
		return NewEffect(convert(x.Expr()), x.Expanded())
	case PositionalArgument:
		return NewPositionalArgument(convert(x.Value()), x.Expanded())
	case Signature:
		ks := make([]OptionalParameter, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, convert(k).(OptionalParameter))
		}

		return NewSignature(x.Positionals(), x.RestPositionals(), ks, x.RestKeywords())
	case Switch:
		cs := make([]SwitchCase, 0, len(x.Cases()))

		for _, c := range x.Cases() {
			cs = append(cs, convert(c).(SwitchCase))
		}

		return NewSwitch(convert(x.Value()), cs, convert(x.DefaultCase()))
	case SwitchCase:
		return NewSwitchCase(x.Pattern(), convert(x.Value()))
	}

	panic(fmt.Errorf("Invalid value: %#v", x))
}
