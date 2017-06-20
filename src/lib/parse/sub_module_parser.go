package parse

type subModuleParser struct {
	state
}

// NewSubModuleParser creates a new main module parser.
func NewSubModuleParser(file, source string) Parser {
	return &subModuleParser{newState(file, source)}
}

// Parse parses a statement.
func (p *subModuleParser) Parse(macros map[string]func(...interface{}) interface{}) (interface{}, error) {
	s, err := p.statement(p.importModule(), p.let())()

	if err != nil {
		return nil, err
	}

	return s, nil
}

// Finished checks if parsing is finished or not.
func (p *subModuleParser) Finished() bool {
	_, err := p.Exhaust(p.blank())()
	return err == nil
}
