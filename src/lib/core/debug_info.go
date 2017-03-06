package core

// DebugInfo represents metadata of a call.
type DebugInfo struct {
	file       string
	lineNumber int
	source     string // source code at lineNumber in file
}

// NewDebugInfo creates a DebugInfo.
func NewDebugInfo(file string, lineNumber int, source string) DebugInfo {
	return DebugInfo{file, lineNumber, source}
}
