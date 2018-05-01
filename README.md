# Cloe

[![Circle CI](https://img.shields.io/circleci/project/github/cloe-lang/cloe.svg?style=flat-square)](https://circleci.com/gh/cloe-lang/cloe)
[![Coveralls](https://img.shields.io/coveralls/cloe-lang/cloe.svg?style=flat-square)](https://coveralls.io/github/cloe-lang/cloe)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloe-lang/cloe?style=flat-square)](https://goreportcard.com/report/github.com/cloe-lang/cloe)
[![License](https://img.shields.io/github/license/cloe-lang/cloe.svg?style=flat-square)](https://opensource.org/licenses/MIT)

<div align="center">
  <img src="https://raw.githubusercontent.com/cloe-lang/icon/master/spaced.png" alt="logo"/>
</div>

Cloe is a functional programming language.
It aims to be simple, canonical, and practical.

## Installation

```
go get github.com/cloe-lang/cloe/src/cmd/cloe
```

Go 1.8+ is required.

## Features

### Implicit parallelism and concurrency

Inherently, programs written in the language run concurrently and can run
parallelly.
You don't need to parallelize your code explicitly as the language runtime does
it all by itself.

You can also increase parallelism of your code or run functions
sequentially using `par` or `seq` primitives.

### Semi-purely functional programming

Reading data outside programs are regarded as pure and side effects going
outside programs are treated as impure.
Impure function calls in pure functions are inspected and raise errors at
runtime.

## Documentation

[Here](https://cloe-lang.org).

## Examples

### Hello, world!

```
(write "Hello, world!")
```

### HTTP server

```
(import "http")

(def (handler request)
  ((request "respond") "Hello, world!"))

(let requests (http.getRequests ":8080"))

..(map handler requests)
```

See [examples](examples) directory for more.

## License

[MIT](LICENSE)
