# Language Specification (by examples)

> Warning: Too many parentheses ahead!

## Basics

### Simple function calls

```
(function argument)
(function argument1 argument2)
(function argument1 argument2 argument3 argument4 argument5)
```

## Data structures

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
```

### Mutually recursive function definition `mr`

TBD

### Macro definition `macro`

TBD

## Built-in functions

TBD
