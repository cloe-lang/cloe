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

		ds := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, dict := range x.ExpandedDicts() {
			ds = append(ds, convert(dict))
		}

		return NewArguments(ps, ks, ds)
	case Import:
		return x
	case KeywordArgument:
		return NewKeywordArgument(x.Name(), convert(x.Value()))
	case LetFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			ls = append(ls, convert(l))
		}

		return NewLetFunction(
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
		fs := make([]LetFunction, 0, len(x.LetFunctions()))

		for _, f := range x.LetFunctions() {
			fs = append(fs, convert(f).(LetFunction))
		}

		return NewMutualRecursion(fs, x.DebugInfo())
	case OptionalArgument:
		return NewOptionalArgument(x.Name(), convert(x.DefaultValue()))
	case Effect:
		return NewEffect(convert(x.Expr()), x.Expanded())
	case PositionalArgument:
		return NewPositionalArgument(convert(x.Value()), x.Expanded())
	case Signature:
		ps := make([]OptionalArgument, 0, len(x.PosOpts()))

		for _, p := range x.PosOpts() {
			ps = append(ps, convert(p).(OptionalArgument))
		}

		ks := make([]OptionalArgument, 0, len(x.KeyOpts()))

		for _, k := range x.KeyOpts() {
			ks = append(ks, convert(k).(OptionalArgument))
		}

		return NewSignature(x.PosReqs(), ps, x.PosRest(), x.KeyReqs(), ks, x.KeyRest())
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
