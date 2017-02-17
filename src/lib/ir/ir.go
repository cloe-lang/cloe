package ir

import "../vm"

type IRThunk struct {
	function interface{}
	args     IRArguments
}

func IRApp(f interface{}, args IRArguments) IRThunk {
	return IRThunk{
		function: f,
		args:     args,
	}
}

func (t IRThunk) compile(args []*vm.Thunk) *vm.Thunk {
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

type IRArguments struct {
	positionals   []IRPositionalArgument
	keywords      []IRKeywordArgument
	expandedDicts []interface{}
}

func NewIRArguments(ps []IRPositionalArgument, ks []IRKeywordArgument, dicts []interface{}) IRArguments {
	return IRArguments{
		positionals:   ps,
		keywords:      ks,
		expandedDicts: dicts,
	}
}

func NewIRPositionalArguments(xs ...interface{}) IRArguments {
	ps := make([]IRPositionalArgument, len(xs))

	for i, x := range xs {
		ps[i] = NewIRPositionalArgument(x, false)
	}

	return NewIRArguments(ps, []IRKeywordArgument{}, []interface{}{})
}

type IRPositionalArgument struct {
	value    interface{}
	expanded bool
}

func NewIRPositionalArgument(ir interface{}, expanded bool) IRPositionalArgument {
	return IRPositionalArgument{
		value:    ir,
		expanded: expanded,
	}
}

func (p IRPositionalArgument) compile(args []*vm.Thunk) vm.PositionalArgument {
	return vm.NewPositionalArgument(compileExpression(args, p.value), p.expanded)
}

type IRKeywordArgument struct {
	name  string
	value interface{}
}

func NewIRKeywordArgument(n string, ir interface{}) IRKeywordArgument {
	return IRKeywordArgument{
		name:  n,
		value: ir,
	}
}

func (k IRKeywordArgument) compile(args []*vm.Thunk) vm.KeywordArgument {
	return vm.NewKeywordArgument(k.name, compileExpression(args, k.value))
}
