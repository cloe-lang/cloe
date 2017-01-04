package desugar

func Desugar(x interface{}) interface{} {
	return newState().module(x)
}

func (s *state) module(x interface{}) interface{} {
	return mapFunc(s.statement, x.([]interface{}))
}

func (s *state) statement(x interface{}) interface{} {
	xs := x.([]interface{})

	if xs[0] == "let" {
		return append(xs[:1], s.expr(xs[2]))
	}

	return s.app(xs)
}

func (s *state) expr(x interface{}) interface{} {
	if app := toApp(x); app != nil {
		return s.app(app)
	}

	return x
}

func (s *state) app(xs []interface{}) interface{} {
	// TODO
	return xs
}

func toApp(x interface{}) []interface{} {
	app, ok := x.([]interface{})

	if !ok {
		return nil
	}

	f := app[0]

	if f == "\\" {
		return nil
	}

	return app
}

func mapFunc(f func(interface{}) interface{}, xs []interface{}) []interface{} {
	ys := make([]interface{}, len(xs))

	for i, x := range xs {
		ys[i] = f(x)
	}

	return ys
}
