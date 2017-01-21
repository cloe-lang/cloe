# Language

- Dynamically typed
- Object == Dictionary
  - Like Clojure and JavaScript
- Support only Ad-hoc polymorphism
  - No subtyping
  - No generics (Dynamic typing covers it.)
- Follow the zen of Python
- 3 special forms
  - `let`: Function or constant definition
  - `rec`: Mutually recursive function definition
  - `macro`: Macro definition


## Examples

```
(let answer 42)
```

```
(let (inc x) (+ x 1))
```

```
(let (factorial n) (if (= n 0) 1 (* n (factorial (- n 1)))))
```

```
; OCaml style mutual recursion
(rec
  (let (even n) (if (= n 0) true  (odd  (- n 1))))
  (let (odd  n) (if (= n 0) false (even (- n 1)))))
```

```
(macro (call func ..args) (func ..args))
```
