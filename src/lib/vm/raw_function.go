package vm

type RawFunction func(Dictionary) Object // Environment -> Object

func (f RawFunction) Call(args Dictionary) Object {
	return f(args)
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
