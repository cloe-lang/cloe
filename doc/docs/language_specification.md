# Language Specification (by examples)

> Warning: Too many parentheses ahead!

## Basics

### Simple function calls

A simple function call is represented by a function and some positional
arguments.

```
(function argument)
(function argument1 argument2)
(function argument1 argument2 argument3 argument4 argument5)
```

## Data types

There are number, string, bool, nil, list, dictionary, and function types.

### Number

```
123
1.1
010
0x10
-123
(+ 1 2 3)
(- 3 2 1)
(* 1 2 3)
(/ 1 2)
(** 2 3)
(mod 42 5)
```

### String

```
"Foo bar baz."
```

### Bool

```
true
false
```

### Nil

```
nil
```

### List

```
[1 2 3]
```

### Dictionary

```
{"key1" 123 "key2" 456}
```

## Function calls

In addition to positional arguments, Tisp supports keyword arguments,
positional rest arguments and keyword rest arguments.

```
(function positionalArgument)
(function positionalArgument1 positionalArgument2)
(function ..expandedList)
(function . keywordArgument itsValue)
(function . keywordArgument1 itsValue1 keywordArgument2 itsValue2)
(function . ..expandedDictionary)
(function positionalArgument1 positionalArgument2 ..expandedList1
          positionalArgument3 ..expandedList2
          .
          keywordArgument1 itsValue1
          keywordArgument2 itsValue2
          ..expandedDictionary1
          ..expandedDictionary2)
```

## Special forms

### Function or variable definition `let`

```
(let myVariable 123)
(let (myFunc x) (+ x 42))
(let (factorial n) (if (= n 0) 1 (* n (factorial (- n 1)))))
```

### Mutually-recursive functions `mr`

```
(mr
  (let (even? n)
       (if (= n 0) true (odd? (- n 1))))
  (let (odd? n)
       (if (= n 0) false (even? (- n 1)))))
```

### Macro definition `macro`

TBD

## Built-in functions

TBD
