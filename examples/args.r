(let kwargs {'y5 123 'y6 456})
(let args '(1 2 3 4 "foo-bar-baz"))

; This is comment!!!

((\ (x1 x2 (x3 123) (x4 456) args.. y1 (y2 123) y3 (y4 456) kwargs...) 42)
 1 2 3 4 args.. y1 123 y3 456 kwargs...)
