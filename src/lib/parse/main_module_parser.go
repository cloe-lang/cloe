package parse

// MainModuleParser is a parser for main modules.
type MainModuleParser struct {
	state
}

// NewMainModuleParser creates a new main module parser.
func NewMainModuleParser(file, source string) Parser {
	return &MainModuleParser{newState(file, source)}
}

// Parse parses a statement.
func (p *MainModuleParser) Parse() (interface{}, error) {
	s, err := p.statement(p.importModule(), p.let(), p.output())()

	if err != nil {
		return nil, err
	}

	return s, nil
}

// Finished checks if parsing is finished or not.
func (p *MainModuleParser) Finished() bool {
	_, err := p.Exhaust(p.blank())()
	return err == nil
}
