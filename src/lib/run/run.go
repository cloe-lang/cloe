package run

import "../core"

func Run(ts []*core.Thunk) {
	for _, t := range ts {
		go t.Eval() // TODO: Ensure results are Outputs.
	}

	for _, t := range ts {
		t.Eval()
	}
}
