package gensym

import (
	"fmt"
	"strings"
	"sync/atomic"
)

var index uint64

// GenSym creates a unique symbol as a string from its components.
func GenSym(ss ...string) string {
	return strings.Join(append([]string{"gensym", fmt.Sprint(atomic.AddUint64(&index, 1) - 1)}, ss...), "$")
}
