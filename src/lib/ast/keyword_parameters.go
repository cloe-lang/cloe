package ast

import "strings"

type keywordParameters struct {
	parameters []OptionalParameter
	rest       string
}

func (ks keywordParameters) names() []string {
	ns := make([]string, 0, len(ks.parameters)+1)

	for _, k := range ks.parameters {
		ns = append(ns, k.name)
	}

	if ks.rest != "" {
		ns = append(ns, ks.rest)
	}

	return ns
}

func (ks keywordParameters) String() string {
	ss := make([]string, 0, len(ks.parameters)+1)

	for _, o := range ks.parameters {
		ss = append(ss, o.String())
	}

	if ks.rest != "" {
		ss = append(ss, ".."+ks.rest)
	}

	return strings.Join(ss, " ")
}
