## Language

- Everything is an object
  - Every value is wrapped with a dictionary.
- Support only Ad-hoc polymorphism
  - No subtyping
  - No generics (dynamic typing covers it.)
  - Interface support
- Follow the zen of Python
- 4 special forms
  - `let`: Constant definition
  - `def`: Recursive or non-recursive function definition
  - `rec`: Mutually recursive function definition
  - `out`: Output definition


## OCaml-style mutual recursion

```
(rec ((even n) (if (= n 0) true  (odd  (- n 1))))
     ((odd  n) (if (= n 0) false (even (- n 1)))))
```
