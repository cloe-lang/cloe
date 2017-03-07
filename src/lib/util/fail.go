package util

import "log"

// Fail makes a program fail printing log message.
func Fail(s string, xs ...interface{}) {
	log.Fatalf(s, xs)
}
