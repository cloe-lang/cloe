package run

import "../vm"

func Run(ts []*vm.Thunk) {
	for _, t := range ts {
		go t.Eval() // TODO: Ensure results are Outputs.
	}

	for _, t := range ts {
		t.Eval()
	}
}
