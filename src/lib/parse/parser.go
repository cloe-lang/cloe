package parse

// Parser is a parser for a module.
type Parser interface {
	Parse() (interface{}, error)
}
