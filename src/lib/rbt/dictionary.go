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

func (d Dictionary) Remove(k interface{}) Dictionary {
	return Dictionary{d.Tree.Remove(k)}
}

func (d Dictionary) FirstRest() (interface{}, interface{}, Dictionary) {
	x, t := d.Tree.FirstRest()

	if x == nil {
		return nil, nil, Dictionary{NewTree(d.less)}
	}

	kv := x.(keyValue)
	return kv.Key, kv.Value, Dictionary{t}
}

func (d Dictionary) Merge(dd Dictionary) Dictionary {
	return Dictionary{d.Tree.Merge(dd.Tree)}
}
