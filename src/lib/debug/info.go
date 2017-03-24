package debug

import (
	"fmt"
	"runtime"
)

// Info represents metadata of a call.
type Info struct {
	file       string
	lineNumber int
	source     string // source code at lineNumber in file
}

// NewInfo creates a Info.
func NewInfo(file string, lineNumber int, source string) Info {
	return Info{file, lineNumber, source}
}

// NewGoInfo creates a Info of debug information about Go source.
func NewGoInfo(skip int) Info {
	_, file, line, ok := runtime.Caller(skip + 1)

	if !ok {
		panic("runtime.Caller failed.")
	}

	return NewInfo(file, line, "")
}

// Lines returns string representation of Info which can be printed on stdout or
// stderr as is.
func (i Info) Lines() string {
	return fmt.Sprintf("%s:%d:\t%s\n", i.file, i.lineNumber, i.source)
}
