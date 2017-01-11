# To do

- `ir <filename>` bin
- `def` desugarer
- multiple desugarers
- Make every `Object` stringable
- Add Keyable interface for String, Nil, Bool, Number
- vm/functions.go with S, K, I, and Y and Ys (mutual recursion)
  - Y and Ys should always be trampolined if possible.
- s/vm/core/g?
- Allow Functions to return Thunks for tail call elimination
  - f : (...*Thunk) -> *Thunk
  - f : (...Object) -> *Thunk
  - NewTOFunction, NewOOFunction, NewTTFunction, NewOTFunction
- Add Thunk.EvalStrict() to evaluate Object(*Thunk)s and always return normal
  Objects
- Trampoline last app like `(\ (x...) (trampoline ...))` in Y
- `Thunk.Eval()` should return `TrampolinedThunk { thunk: *Thunk, trampoliner: *Thunk }`
