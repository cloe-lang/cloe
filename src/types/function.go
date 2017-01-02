package types

type RawFunction func(Dictionary) Object

func (f RawFunction) Call(args Dictionary) Object {
	return f(args)
}
