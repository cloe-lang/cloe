package rbt

type keyValue struct {
	Key, Value interface{}
}

type Dictionary struct{ Tree }

func NewDictionary(less func(interface{}, interface{}) bool) Dictionary {
	lessKV := func(x1, x2 interface{}) bool {
		if kv, ok := x1.(keyValue); ok {
			x1 = kv.Key
		}

		if kv, ok := x2.(keyValue); ok {
			x2 = kv.Key
		}

		return less(x1, x2)
	}

	return Dictionary{NewTree(lessKV)}
}

func (d Dictionary) Insert(k, v interface{}) Dictionary {
	return Dictionary{d.Tree.Insert(keyValue{k, v})}
}

func (d Dictionary) Search(k interface{}) (interface{}, bool) {
	kv, ok := d.Tree.Search(k)

	if !ok {
		return nil, false
	}

	return kv.(keyValue).Value, ok
}

func (d Dictionary) Remove(k interface{}) (Dictionary, bool) {
	t, ok := d.Tree.Remove(k)
	return Dictionary{t}, ok
}

type FirstRestKVFunc func() (interface{}, interface{}, FirstRestKVFunc)

func (d Dictionary) FirstRest() (interface{}, interface{}, FirstRestKVFunc) {
	return convertFunc(d.Tree.FirstRest)()
}

func convertFunc(f FirstRestFunc) FirstRestKVFunc {
	return func() (interface{}, interface{}, FirstRestKVFunc) {
		x, f := f()

		if x == nil {
			return nil, nil, nil
		}

		kv := x.(keyValue)
		return kv.Key, kv.Value, convertFunc(f)
	}
}
