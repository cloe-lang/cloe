# Notes

## Closure

```
(let foo 123)
(print (map (\ (x) (+ x foo)) list))
```

is equivalent to:

```
(let foo 123)
(let f (\ (x y) (+ y x)))
(print (map (partial f foo) list))
```

### Nondeterminism module

```
(let (rally ..objects)
     "Objects -> List. Sort by time when they are evaluated."
     (...))
```

## Keeping purity of functions

Every output function should return `OutputError`. And, the error is ignored by
`sync` (`cause`?) and top-level call of `Thunk.Eval()`.

```
outputFunction : X1 -> X2 -> ... -> Xn -> Error (OutputError)
```

## `rally` vs `outs`

```
rally : [a] -> [a]
outs : [Output] -> Output
```

## Named types like Go language?

It may inhibit potentially polymorphic code.
(e.g. `(isType? x "POSRecord")` instead of `(haveKey x "price")`)
