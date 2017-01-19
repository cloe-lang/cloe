# Language

- Object == Dictionary
  - Like Clojure and JavaScript
- Support only Ad-hoc polymorphism
  - No subtyping
  - No generics (Dynamic typing covers it.)
- Follow the zen of Python
- 4 special forms
  - `\`: Lambda expression
  - `def`: Non-recursive function or constant definition
  - `defr`: (Mutually) recursive function definition
  - `macro`: Macro definition


## Examples

```
(\ (x) (+ x 1))
```

```
(def answer 42)
```

```
(def (inc x) (+ x 1))
```

```
(defr (factorial n) (if (= n 0) 1 (* n (factorial (- n 1)))))
```

```
; OCaml style mutual recursion
(defr
  (even n) (if (= n 0) true  (odd  (- n 1)))
  (odd  n) (if (= n 0) false (even (- n 1))))
```

```
(macro (call func ..args) (func ..args))
```
