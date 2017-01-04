package vm

type Array []Object

func (a Array) ToObjectSlice() []Object {
	return ([]Object)(a)
}
