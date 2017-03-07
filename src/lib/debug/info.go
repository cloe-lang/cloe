package debug

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
