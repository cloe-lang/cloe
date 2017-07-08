package match

type patternType int

const (
	listPattern patternType = iota
	dictPattern
	namePattern
)
