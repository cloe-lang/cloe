package core

import "fmt"

func sprint(x interface{}) StringType {
	return StringType(fmt.Sprint(x))
}
