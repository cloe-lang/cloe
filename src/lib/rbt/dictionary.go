package rbt

type keyValue struct {
	Key, Value interface{}
}

// Dictionary represents a dictionary which contains keys and corresponding values.
type Dictionary struct{ Tree }

// NewDictionary creates an empty dictionary.
func NewDictionary(compare func(interface{}, interface{}) int) Dictionary {
	compareKV := func(x1, x2 interface{}) int {
		if kv, ok := x1.(keyValue); ok {
			x1 = kv.Key
		}

		if kv, ok := x2.(keyValue); ok {
			x2 = kv.Key
		}

		return compare(x1, x2)
	}

	return Dictionary{NewTree(compareKV)}
}

// Insert inserts a key and a corresponding value to a dictionary.
func (d Dictionary) Insert(k, v interface{}) Dictionary {
	return Dictionary{d.Tree.Insert(keyValue{k, v})}
}

// Search searches a key in dictionary, and returns true if it is found or false otherwise.
func (d Dictionary) Search(k interface{}) (interface{}, bool) {
	kv, ok := d.Tree.Search(k)

	if !ok {
		return nil, false
	}

	return kv.(keyValue).Value, ok
}

// Remove removes a key from a dictionary.
func (d Dictionary) Remove(k interface{}) Dictionary {
	return Dictionary{d.Tree.Remove(k)}
}

// FirstRest returns a value which was in the original dictionary and a new
// dictionary without it.
func (d Dictionary) FirstRest() (interface{}, interface{}, Dictionary) {
	x, t := d.Tree.FirstRest()

	if x == nil {
		return nil, nil, Dictionary{NewTree(d.compare)}
	}

	kv := x.(keyValue)
	return kv.Key, kv.Value, Dictionary{t}
}

// Merge merges 2 dictionaries.
// If there are duplicate keys between them, the second dictionary's values are used.
func (d Dictionary) Merge(dd Dictionary) Dictionary {
	return Dictionary{d.Tree.Merge(dd.Tree)}
}
