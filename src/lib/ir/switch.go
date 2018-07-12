package ir

import "github.com/cloe-lang/cloe/src/lib/core"

type switchData struct {
	Cases       []caseData
	DefaultCase int
}

func newSwitchData(cs []caseData, dc int) switchData {
	return switchData{cs, dc}
}

type caseData struct {
	Value         core.Value
	VariableIndex int
}

func newCaseData(v core.Value, i int) caseData {
	return caseData{v, i}
}
