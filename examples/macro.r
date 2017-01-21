(macro (unless b t f) `(if b f t))

(macro (foo) 123)
(print (+ (foo) 42))

; Same as above
(macro (foo) `123)
(print (+ (foo) bar))

; Invalid and useless too; just use let.
(macro foo 123)
