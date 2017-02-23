; . after ..args is not necessary. But it can be forced for readability.
(let (func x1 x2 (x3 123) (x4 456) ..args . y1 y2 (y3 123) (y4 456) ..kwargs)
     (+ x1 x2 x3 x4))

(func 1 2 3 ..[1 "foo" "bar"] . y1 123 y3 456 foo 2049 ..{"y4" 123 "y6" 456})
