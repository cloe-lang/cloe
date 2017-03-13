# Tisp

[![wercker status](https://app.wercker.com/status/68b66e4881f08974e109a864520e369f/s/master "wercker status")](https://app.wercker.com/project/byKey/68b66e4881f08974e109a864520e369f)
[![codecov](https://codecov.io/gh/raviqqe/tisp/branch/master/graph/badge.svg)](https://codecov.io/gh/raviqqe/tisp)
[![codebeat badge](https://codebeat.co/badges/3a45a98c-ad7d-4a0a-8011-241f0ae4682c)](https://codebeat.co/projects/github-com-raviqqe-tisp-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/raviqqe/tisp)](https://goreportcard.com/report/github.com/raviqqe/tisp)
[![License: MIT](https://img.shields.io/badge/license-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![inputs to outputs](img/inputs_to_outputs.png)

Tisp is a "Time is Space" programming language.
It aims to be simple, canonical, and practical.

## Features

- Purely functional programming
  - Impure function calls in pure functions are detected and raise errors at
    runtime.
- Implicit parallelism and concurrency
  - Most of the time, you don't need to parallelize your code explicitly.
    Tisp does it all for you!
  - Inherently, codes written in Tisp can be run concurrently.
- Optional injection of parallelism and causality
  - You can also increase parallelism of your code or run functions
    sequentially using `par` or `seq` primitives.
- Asynchronous IO
- Dynamic typing

## Current state

See [the issues](https://github.com/raviqqe/tisp/issues).
Tisp is actively developed and work in progress.
Stay tuned!

## FAQ

### What languages is Tisp inspired by?

The following is their list and ideas which come from them.

- Haskell
  - The concept of "Time is Space"
  - Lazy evaluation and data structures to realize it
- Clojure
  - Everything is a value; no object system
  - The syntax and macro system
- OCaml
  - The syntax which is close to pure lambda calculus and mutual recursion

## License

[MIT](LICENSE)
