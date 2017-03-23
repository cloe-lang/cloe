# Tisp

The "Time is Space" programming language

Tisp is a functional programming language with implicit parallelism and
concurrency.
It aims to be simple, canonical, and practical.

## TL;DR

- Tisp evaluates every output of a program (such as printing a string and
  sending an HTTP response to a client) parallelly and concurrently by default
  leaving synchronization optional.
- Tisp keeps the other part of a program pure (i.e. without any output).
- Therefore, every program in Tisp can run parallelly and concurrently with
  nothing special!

## Background

Existing programming languages are synchronous by default.
For example, when you write code in C like the below, the statements are run
sequentially one by one even if they can be run in parallel.

```c
send(http_response);
send(another_http_response);
```

Therefore, many libraries, frameworks, and language features for parallel,
concurrent, and asynchronous programming have emerged in recent years.

Tisp takes the different way to deal with this problem.
It lets you write programs which are asynchronous, concurrent, and parallel
inherently.

## How?

> Warning: This is an immature concept.
> Please post any feedback to improve or fix it in
> [the Google group](https://groups.google.com/forum/#!forum/tisp-aliens).

Tisp assumes that *time is space*.
If *time is space*, everything is constant because nothing changes over time.
This fact lets us write programs as functions which map their inputs to
outputs.

Let's think about a program which takes input data `input[t]` and outputs
something `output[t]` at every time step `t`.
Then, every `output[t]` should be calculated from all inputs in the past.

```
input[0], input[1], ..., input[t-1] -> output[t]
```

This is true for any `t`.
Therefore, a program can be represented as a function which maps its inputs to
outputs.

```
program : input[0], ..., input[T] -> output[0], ..., output[T]
```

where `T` is infinite generally.
Then, we can extract great parallelism of a program because the program as a
function represents just data dependency.

Tisp evaluates the outputs concurrently and parallelly.
Programmers must specify their causality using the `seq` primitive function if
necessary.
