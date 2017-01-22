package ir

func Run(instrs []interface{}) {
	os := newCompiler().compile(instrs)

	for _, o := range os {
		go func() {
			// TODO: Ensure results are Outputs.
			o.Eval()
		}()
	}

	for _, o := range os {
		o.Eval()
	}
}
