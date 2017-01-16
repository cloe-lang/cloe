# VM

## Concepts

- Object oriented
  - Primitive types and methods to manipulate them

## Types

- Number
  - DEC64
- String
  - Compressed
- List
- Dictionary
  - Persistent data structure
- Nil
  - is Evil
- Function
  - [Thunk] -> Thunk
- Closure
  - May remove Function
- Set?

## Roadmap

- Tier 1
  - Pure lambda calculus
  - `sync` primitive
- Tier 2
  - IO
- Tier 3
  - Rally sort
  - DEC64
- Tier ?
  - Persistent z(t) support
    - Save objects in files
    - May not necessary and covered by `sync` primitive

## Instructions

- + :
  - Number -> [Number] -> Number
  - String -> [String] -> Number
- - : Number -> [Number] -> Number
- * : Number -> [Number] -> Number
- / : Number -> [Number] -> Number
- // : Number -> [Number] -> Number (floor division)
- ^ : Number -> [Number] -> Number
