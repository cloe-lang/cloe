package ir

import "fmt"

type variables struct {
	nameToIndex  map[string]int
	numVariables int
}

func newVariables(m map[string]int, n int) variables {
	return variables{m, n}
}

func (vs *variables) AddNamedVariable(s string) int {
	i := vs.AddVariable()
	vs.BindNamedVariable(s, i)
	return i
}

func (vs *variables) AddVariable() int {
	vs.numVariables++
	return vs.numVariables - 1
}

func (vs *variables) BindNamedVariable(s string, i int) {
	vs.nameToIndex[s] = i
}

func (vs *variables) GetNamedVariable(s string) (int, error) {
	i, ok := vs.nameToIndex[s]

	if !ok {
		return 0, fmt.Errorf("name, %s not found", s)
	}

	return i, nil
}
