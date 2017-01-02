# Notes

## Mapping of arguments to names

Can be just sugar.

```
(let list '(123 456 789))

(let func (\ (*x) x))

(= (func 123 456 *list 789)
   (func * x (concat '(123 456) list '(789))))
```
