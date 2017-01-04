# VM

## Concepts

- Object oriented
  - Primitive types and methods to manipulate them

## Types

- Number (N)
  - Numbers in an Array construct a String.
  - DEC64
- Array (A)
  - Used to construct Lists
  - Copy for each operation to realize immutability
- Dictionary (D)
  - Used to construct Sets
  - Persistent data structure
- None
  - Evil
- Function
  - Dictionary -> Object
  - `(\x x)` takes `{ "x": Object }` and returns `[x]`
  - `(\x *x)` takes `{ "x": List-like }` and returns `x`

## Roadmap

- Tier 1
  - Pure lambda calculus
  - `out` operator
  - IEEE 754
- Tier 2
  - Persistent z(t) support
    - May not necessary
  - Mutual recursion support
  - DEC64

## Instructions

- + : Number -> Number -> Number
- - : Number -> Number -> Number
- * : Number -> Number -> Number
- / : Number -> Number -> Number
- // : Number -> Number -> Number (floor division)
- ^ : Number -> Number -> Number
