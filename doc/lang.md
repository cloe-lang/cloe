## Language

- Everything is an object
  - Every value is wrapped with a dictionary.
- Support only Adhoc polymorphism
  - No subtyping
  - No generics (dynamic typing covers it.)
  - Interface support

# OCaml style mutual recursion

```
(rec
  (def (even n) (if (= n 0) true  (odd (- n 1))))
  (def (odd  n) (if (= n 0) false (even (- n 1)))))
```
