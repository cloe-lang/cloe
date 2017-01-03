package parse

func toString(any interface{}) string {
	xs := any.([]interface{})
	rs := make([]rune, len(xs))

	for i, x := range xs {
		rs[i] = x.(rune)
	}

	return string(rs)
}
