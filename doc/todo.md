# To do

- `ir <filename>` bin
- `def` desugarer
- multiple desugarers

- Pattern match
- Type system?
  - May be gradual one
- Slab allocator for Thunks?

## vm

- Split out `concat` function of `stringType` from `add` generic function.
  - Then, `add` and `mul` for `numberType` can take no argument.
- `listable` implementation for `stringType`

- s/vm/core/g?
