package desugar

func Desugar(x interface{}) interface{} {
	return newState().desugar(x)
}

func (s *state) desugar(x interface{}) interface{} {
	if xs, ok := x.([]interface{}); ok && isApp(xs) {
		return xs
	} else if ok {
		return mapFunc(s.desugar, xs)
	}

	return x
}

func isApp(xs []interface{}) bool {
	head := xs[0]
	return head != "\\" && head != "let" && head != "quote"
}

func mapFunc(f func(interface{}) interface{}, xs []interface{}) []interface{} {
	ys := make([]interface{}, len(xs))

	for i, x := range xs {
		ys[i] = f(x)
	}

	return ys
}
