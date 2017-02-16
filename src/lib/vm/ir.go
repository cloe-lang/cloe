package vm

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

func (t IRThunk) compile(args []*Thunk) *Thunk {
	ps := make([]PositionalArgument, len(t.args.positionals))

	for i, p := range t.args.positionals {
		ps[i] = p.compile(args)
	}

	ks := make([]KeywordArgument, len(t.args.keywords))

	for i, k := range t.args.keywords {
		ks[i] = k.compile(args)
	}

	ds := make([]*Thunk, len(t.args.expandedDicts))

	for i, d := range t.args.expandedDicts {
		ds[i] = compileExpression(args, d)
	}

	return App(compileExpression(args, t.function), NewArguments(ps, ks, ds))
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

func (p IRPositionalArgument) compile(args []*Thunk) PositionalArgument {
	return NewPositionalArgument(compileExpression(args, p.value), p.expanded)
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

func (k IRKeywordArgument) compile(args []*Thunk) KeywordArgument {
	return NewKeywordArgument(k.name, compileExpression(args, k.value))
}
