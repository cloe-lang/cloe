package compile

import "../vm"

type Gettable interface {
	get(string) *vm.Thunk
}
