package util

// StringsToAnys converts []string into []interface{}.
func StringsToAnys(ss []string) []interface{} {
	xs := make([]interface{}, 0, len(ss))

	for _, s := range ss {
		xs = append(xs, s)
	}

	return xs
}
