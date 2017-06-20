package parse

// Parser is a parser for a module.
type Parser interface {
	Parse(macros map[string]func(...interface{}) interface{}) (interface{}, error)
	Finished() bool
}
