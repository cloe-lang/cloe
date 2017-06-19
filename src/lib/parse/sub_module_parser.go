package parse

// SubModuleParser is a parser for main modules.
type SubModuleParser struct {
	state
}

// NewSubModuleParser creates a new main module parser.
func NewSubModuleParser(file, source string) Parser {
	return &SubModuleParser{newState(file, source)}
}

// Parse parses a statement.
func (p *SubModuleParser) Parse() (interface{}, error) {
	s, err := p.statement(p.importModule(), p.let())()

	if err != nil {
		return nil, err
	}

	return s, nil
}

// Finished checks if parsing is finished or not.
func (p *SubModuleParser) Finished() bool {
	_, err := p.Exhaust(p.blank())()
	return err == nil
}
