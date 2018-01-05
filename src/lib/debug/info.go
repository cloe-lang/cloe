package debug

import (
	"fmt"
	"runtime"
)

// Info represents metadata of a call.
type Info struct {
	file                     string
	lineNumber, linePosition int
	source                   string // source code at lineNumber in file
}

// NewInfo creates a Info.
func NewInfo(file string, lineNumber, linePosition int, source string) Info {
	return Info{file, lineNumber, linePosition, source}
}

// NewGoInfo creates a Info of debug information about Go source.
func NewGoInfo(skip int) Info {
	if !Debug {
		return Info{}
	}

	_, file, line, ok := runtime.Caller(skip + 1)

	if !ok {
		panic("runtime.Caller failed.")
	}

	return NewInfo(file, line, -1, "")
}

// Lines returns string representation of Info which can be printed on stdout or
// stderr as is.
func (i Info) Lines() string {
	if i == (Info{}) {
		return ""
	}

	p := "NA"

	if i.linePosition > 0 {
		p = fmt.Sprint(i.linePosition)
	}

	return fmt.Sprintf("%s:%d:%s:\t%s\n", i.file, i.lineNumber, p, i.source)
}
