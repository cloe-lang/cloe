package run

import "github.com/raviqqe/tisp/src/lib/core"

func Run(ts []*core.Thunk) {
	for _, t := range ts {
		go t.Eval() // TODO: Ensure results are Outputs.
	}

	for _, t := range ts {
		t.Eval()
	}
}
