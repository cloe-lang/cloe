package vm

type Equalable interface {
	Equal(t1, t2 *Thunk) *Thunk
}
