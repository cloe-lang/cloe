(let foo 123)
(let bar 456)

(let (func x y)
  "This function calculate (x + y)^3 + (x + y)^2 + (x + y)^1"
  "It should be very useful."
  (let z (+ x y))
  (+ (^ z 3) (^ z 2) z))

(print (+ foo bar))
