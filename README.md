# Tisp

[![wercker](https://img.shields.io/wercker/ci/wercker/docs.svg?style=flat-square)](https://app.wercker.com/tisp-lang/tisp/runs)
[![codeclimate](https://img.shields.io/codeclimate/github/kabisaict/flow.svg?style=flat-square)](https://codeclimate.com/github/tisp-lang/tisp)
[![Go Report Card](https://goreportcard.com/badge/github.com/tisp-lang/tisp?style=flat-square)](https://goreportcard.com/report/github.com/tisp-lang/tisp)
[![codecov](https://img.shields.io/codecov/c/github/tisp-lang/tisp.svg?style=flat-square)](https://codecov.io/gh/tisp-lang/tisp)
[![MIT license](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Google group](https://img.shields.io/badge/join-us-ff69b4.svg?style=flat-square)](https://groups.google.com/forum/#!forum/tisp-aliens)

![logo](https://raw.githubusercontent.com/tisp-lang/icon/master/landscape.png)

Tisp is a "Time is Space" programming language.
It aims to be simple, canonical, and practical.

This baby project is under heavy development.
Any contributions would be appreciated.
[Join the Google group today!](https://groups.google.com/forum/#!forum/tisp-aliens)

## Current status

See [the issues](https://github.com/tisp-lang/tisp/issues).
Tisp is actively developed and work in progress.
Stay tuned!

## Installation

```
go get github.com/tisp-lang/tisp/src/cmd/tisp
```

You need Go 1.8+.

## Features

- Purely functional programming
  - Impure function calls in pure functions are detected and raise errors at
    runtime.
- Implicit parallelism and concurrency
  - Most of the time, you don't need to parallelize your code explicitly.
    Tisp does it all for you!
  - Inherently, programs written in Tisp run concurrently and can run
    parallelly.
- Optional injection of parallelism and causality
  - You can also increase parallelism of your code or run functions
    sequentially using `par` or `seq` primitives.
    (`par` is not implemented yet.)
- Asynchronous IO
  - Every IO can be run asynchronously by the `par` primitive.
- Dynamic typing

## Documentation

Visit [here](https://tisp-lang.github.io/tisp/).

## Examples

Try scripts in [test directory](test) (`test/*.tisp`) by running:

```
tisp test/$filename.tisp
```

Some examples in [examples directory](examples) work but not all because
some features are not implemented yet.

## Contributing

Please send pull requests, issues, questions or anything else.
See also [the documentation](https://tisp-lang.github.io/tisp/for_developers/)
on how to develop Tisp.

## FAQ

### What languages is Tisp inspired by?

The following is their list with ideas and technologies which come from them.

- Haskell
  - The concept of "Time is Space"
  - Lazy evaluation and data structures to realize it
- Clojure
  - Everything is a value; no object system
  - The syntax and macro system
- OCaml
  - The syntax close to pure lambda calculus and of mutual recursion
- Python
  - The Zen (See `python -c 'import this'`.)
  - The syntax of function calls with positional and keyword arguments
- Go
  - Simplicity
  - Tisp utilizes Go's coroutine runtime.

## License

[MIT](LICENSE)
