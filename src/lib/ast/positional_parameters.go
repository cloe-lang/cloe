package ast

import "strings"

type positionalParameters struct {
	parameters []string
	rest       string
}

func (ps positionalParameters) names() []string {
	ns := ps.parameters

	if ps.rest != "" {
		ns = append(ns, ps.rest)
	}

	return ns
}

func (ps positionalParameters) String() string {
	ss := make([]string, 0, len(ps.parameters)+1)

	ss = append(ss, ps.parameters...)

	if ps.rest != "" {
		ss = append(ss, ".."+ps.rest)
	}

	return strings.Join(ss, " ")
}
