package gensym

import (
	"fmt"
	"sync/atomic"
)

var index uint64

// GenSym creates a unique symbol as a string from its components.
func GenSym() string {
	return "$gensym$" + fmt.Sprint(atomic.AddUint64(&index, 1)-1)
}
