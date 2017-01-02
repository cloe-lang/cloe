(let kwargs {'y5 123 'y6 456})

((\ (x1 x2 (x3 123) (x4 456) *args y1 (y2 123) y3 (y4 456) **kwargs) x)
 1 2 3 4 *list * y1 123 y3 456 **kwargs)

; We don't have to support neither `*[123 456]` nor `**{'y5 123 'y6 456}`
; because they can just be expanded into arguments directly.
; (e.g. `(func 123 456 * y5 123 y6 456)`)
