package vm

type Function struct {
	signature List
	env       Dictionary
	function  func(Dictionary) *Thunk // Environment -> Result
}

func (f Function) Call(args List) *Thunk {
	return f.function(mapArgs(f.env, f.signature, args))
}

func mapArgs(env Dictionary, sig, args List) Dictionary {
	// TODO
	return NewDictionary().(Dictionary)
}

// func CompileFunction(o Object) (RawFunction, error) {
// 	os := o.(List).Slice()

// 	if !len(os) != 3 {
// 		return nil, Error(
// 			"Invalid number of elements in a list representing a function. %#v", os)
// 	}

// 	args := os[1]
// 	body := os[2]

// 	return func(env Dictionary) Object {
// 		if v, ok := Dictioanary.Get(); ok {

// 		} else {
// 		}

// 		return
// 	}
// }
