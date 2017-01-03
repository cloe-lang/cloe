package types

import "github.com/mediocregopher/seq"

type Dictionary seq.HashMap

func NewDictionary(args Dictionary) Object {
	o, ok := args.Get(NewString("key_value_pairs"))

	if !ok {
		panic(`An argument, "key_value_pairs" not found.`)
	}

	a, ok := o.(Array)

	if !ok {
		panic(`Invalid type for an argument "key_value_pairs".`)
	}

	if len(a)%2 != 0 {
		panic("Number of arguments is not even.")
	}

	h := seq.NewHashMap()

	for i := 0; i < len(a); i += 2 {
		h.Set(a[i], a[i+1])
	}

	return (*Dictionary)(h)
}

func (d *Dictionary) Get(k Object) (Object, bool) {
	return (*seq.HashMap)(d).Get(k)
}
