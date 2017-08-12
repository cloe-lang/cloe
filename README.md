# Tisp

[![wercker](https://img.shields.io/wercker/ci/wercker/docs.svg?style=flat-square)](https://app.wercker.com/tisp-lang/tisp/runs)
[![codeclimate](https://img.shields.io/codeclimate/github/kabisaict/flow.svg?style=flat-square)](https://codeclimate.com/github/tisp-lang/tisp)
[![Go Report Card](https://goreportcard.com/badge/github.com/tisp-lang/tisp?style=flat-square)](https://goreportcard.com/report/github.com/tisp-lang/tisp)
[![codecov](https://img.shields.io/codecov/c/github/tisp-lang/tisp.svg?style=flat-square)](https://codecov.io/gh/tisp-lang/tisp)
[![MIT license](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Google group](https://img.shields.io/badge/join-us-ff69b4.svg?style=flat-square)](https://groups.google.com/forum/#!forum/tisp-aliens)

<div align="center">
  <img src="https://raw.githubusercontent.com/tisp-lang/icon/master/shadowed.png" alt="logo"/>
</div>

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
- Asynchronous IO
- Dynamic typing

## Documentation

[Here](https://tisp-lang.gitbooks.io/tisp-programming-language/).

## Examples

Read and try [examples](examples) written in [Gherkin](https://cucumber.io/docs/reference).

## Contributing

Please send pull requests, issues, questions or anything else.
See also [the documentation](https://tisp-lang.gitbooks.io/tisp-programming-language/contribution_guide.html)
on how to develop Tisp.

## License

[MIT](LICENSE)
