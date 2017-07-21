package compile

import "strings"

var builtins = func() environment {
	e := prelude.copy()

	for k, v := range subModule(prelude, "<builtins>", `
	(let (List ..xs) xs)
	`) {
		k = strings.ToLower(k[:1]) + k[1:]
		e.set(k, v)
		e.set("$"+k, v)
	}

	return e
}()
