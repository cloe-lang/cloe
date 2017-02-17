package compile

import "../vm"

type Thunk struct {
	function interface{}
	args     Arguments
}

func App(f interface{}, args Arguments) Thunk {
	return Thunk{
		function: f,
		args:     args,
	}
}

func (t Thunk) compile(args []*vm.Thunk) *vm.Thunk {
	ps := make([]vm.PositionalArgument, len(t.args.positionals))

	for i, p := range t.args.positionals {
		ps[i] = p.compile(args)
	}

	ks := make([]vm.KeywordArgument, len(t.args.keywords))

	for i, k := range t.args.keywords {
		ks[i] = k.compile(args)
	}

	ds := make([]*vm.Thunk, len(t.args.expandedDicts))

	for i, d := range t.args.expandedDicts {
		ds[i] = compileExpression(args, d)
	}

	return vm.App(compileExpression(args, t.function), vm.NewArguments(ps, ks, ds))
}

type Arguments struct {
	positionals   []PositionalArgument
	keywords      []KeywordArgument
	expandedDicts []interface{}
}

func NewArguments(ps []PositionalArgument, ks []KeywordArgument, dicts []interface{}) Arguments {
	return Arguments{
		positionals:   ps,
		keywords:      ks,
		expandedDicts: dicts,
	}
}

func NewPositionalArguments(xs ...interface{}) Arguments {
	ps := make([]PositionalArgument, len(xs))

	for i, x := range xs {
		ps[i] = NewPositionalArgument(x, false)
	}

	return NewArguments(ps, []KeywordArgument{}, []interface{}{})
}

type PositionalArgument struct {
	value    interface{}
	expanded bool
}

func NewPositionalArgument(ir interface{}, expanded bool) PositionalArgument {
	return PositionalArgument{
		value:    ir,
		expanded: expanded,
	}
}

func (p PositionalArgument) compile(args []*vm.Thunk) vm.PositionalArgument {
	return vm.NewPositionalArgument(compileExpression(args, p.value), p.expanded)
}

type KeywordArgument struct {
	name  string
	value interface{}
}

func NewKeywordArgument(n string, ir interface{}) KeywordArgument {
	return KeywordArgument{
		name:  n,
		value: ir,
	}
}

func (k KeywordArgument) compile(args []*vm.Thunk) vm.KeywordArgument {
	return vm.NewKeywordArgument(k.name, compileExpression(args, k.value))
}
