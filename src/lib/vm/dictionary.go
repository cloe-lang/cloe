package vm

import "github.com/mediocregopher/seq"

type Dictionary struct{ *seq.HashMap }

func NewDictionary(os ...Object) Object {
	d := Dictionary{seq.NewHashMap()}

	if len(os)%2 != 0 {
		return NewError("Number of arguments is not even.")
	}

	for i := 0; i < len(os); i += 2 {
		d.Set(os[i], os[i+1])
	}

	return d
}
