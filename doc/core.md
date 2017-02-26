# VM

## Types

- Bool
  - if
- Nil
  - The evil
- Number
  - DEC64
  - +, -, *, /, //, mod, **
  - Bit operators?
- String
  - encode (and decode)

- List
  - (index)
- Dictionary
  - (insert), (index)
- Set
  - (insert)

- Callable
  - Function, Closure
  - No method

- Error

- Array?


## Interfaces

- Exported
  - =
    - Except for Error
  - <, >
    - Number, String, List, Set
  - str
    - All types
  - len, include
    - String, List, Dictionary, Set
  - delete
    - Dictionary, Set (List?)
  - merge
    - String, List, Set, Dictionary
- Internal
  - ordered
    - Except for Error
