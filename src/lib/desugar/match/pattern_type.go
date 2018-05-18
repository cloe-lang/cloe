package match

type patternType int

const (
	listPattern patternType = iota
	dictionaryPattern
	scalarPattern
	namePattern
)
