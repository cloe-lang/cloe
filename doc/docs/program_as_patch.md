# Program As Patch

This page tries to explain Tisp without the concept of "time is space".

## Huh?

In Tisp, we view a program as a set of patches (or diffs) which are applied to
*the world* by its runtime.
Every top-level expressions are patches.
For example, a following program in Tisp represents one consisting of a patch
which changes a state of your terminal (part of *the world*) printing
"Hello, world!" on it.

```
(write "Hello, world!")
```

Another program below is a set of 2 patches.
One prints "Hello, world!" on your terminal.
The other prints "Hello, John!" there.

```
(write "Hello, world!")
(write "Hello, John!")
```

Note that Tisp doesn't care about orders of patches by default.
So, the result of the last program can be either:

```
Hello, world!
Hello, John!
```

or

```
Hello, John!
Hello, world!
```

When you want to apply patches sequentially one by one,
you can use `seq` primitive function to do so.
The result of `seq` function is also a patch.

```
(seq
  (write "Hello, world!")
  (write "Hello, John!"))
```

Definitions of patches are allowed only as top-level expressions or arguments
to the `seq` function.

## Inputs and outputs

Patches are outputs of programs.
It is not enough; we need a way to define inputs of programs.

That said, we do not have to distinguish inputs of programs from usual
expressions in Tisp.
For example, `read` function reads stdin and returns its result as a string.

```
(write (read))
```

The program above just reads stdin and writes it back to stdout.
A caveat is that every expression in Tisp is evaluated lazily; every expression
is evaluated only when it is needed to calculate its result.
For example, the following program reads nothing from stdin and writes "foo" on
terminal.

```
(let (f x) "foo")
(write (f (read)))
```

This is a kind of a dangerous way to define inputs of programs because they are
often patches too.
For example, if a program reads stdin twice, the second result will be an empty
string.
Therefore, programmers must be sure that stdin is read only at one place in
each program in Tisp.

However, with this paradigm, we can write patches as usual expressions and
extract parallelism and concurrency of programs easily.

## Implicit parallelism and concurrency

Tisp runs each patches in a program parallelly and concurrently.
This is how Tisp realizes implicit parallelism and concurrency.
