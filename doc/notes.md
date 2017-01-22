# Notes

## Mapping of arguments to names

Can be just sugar because every variadic function has a certain number of
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

Expressions containing `import` calls must be evaluated at compile time and
a compiler must prevent use of `m` as a dictionary.

Or, treat `import` as just a another special form.

```
(import "directory/module_name") ; statically used by compiler
(print (module_name.calcAnswer 42))
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

## Error handling

- Stick with returning errors
  - Catching exceptions with Go's `defer` is problematic because they can be
    catched only by functions which evaluates erroneous thunks.
- Instantiate VM instructions with debug information at each call in each
  function?

## Do we need `seq` and `par` as Parallel Haskell?

- The former, `seq` is pretty similar to `cause` while `cause` is only for IO
  synchronization.
