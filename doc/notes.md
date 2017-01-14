# Notes

## Mapping of arguments to names

Can be just sugar because every variadic function has certain number of
arguments.

```
(let list '(123 456 789))

(let func (\ (*x) x))

(= (func 123 456 *list 789)
   (func * x (concat '(123 456) list '(789))))
```

## Evaluation steps

1. Parse source code: `string -> []interface{}`
2. Eval source code: `[]interface{} -> *Thunk`

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

## Module system

```
(let m (import "directory/module_name"))
(print (m.calcAnswer 42))
```

### Nondeterminism module

```
(def nd.rally "Objects -> List. Sort by time when they are evaluated."
  (objects...) (...))
```
