package comb

// Parser is a type of parsers as a function which returns a parsing result or
// an error.
type Parser func() (interface{}, error)
