(let answer 42)

(let (inc x) (+ x 1))

(let (factorial n) (if (= n 0) 1 (* n (factorial (- n 1)))))

; OCaml style mutual recursion
(rec
  (let (even n) (if (= n 0) true  (odd  (- n 1))))
  (let (odd  n) (if (= n 0) false (even (- n 1)))))

(macro (call func ..args) (func ..args))
