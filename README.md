# Coel

[![Circle CI](https://img.shields.io/circleci/project/github/coel-lang/coel.svg?style=flat-square)](https://circleci.com/gh/coel-lang/coel)
[![Coveralls](https://img.shields.io/coveralls/coel-lang/coel.svg?style=flat-square)](https://coveralls.io/github/coel-lang/coel)
[![Go Report Card](https://goreportcard.com/badge/github.com/coel-lang/coel?style=flat-square)](https://goreportcard.com/report/github.com/coel-lang/coel)
[![codeclimate](https://img.shields.io/codeclimate/github/kabisaict/flow.svg?style=flat-square)](https://codeclimate.com/github/coel-lang/coel)
[![License](https://img.shields.io/github/license/coel-lang/coel.svg?style=flat-square)](https://opensource.org/licenses/MIT)

<div align="center">
  <img src="https://raw.githubusercontent.com/coel-lang/icon/master/spaced.png" alt="logo"/>
</div>

Coel is a functional programming language.
It aims to be simple, canonical, and practical.

## Installation

```
go get github.com/coel-lang/coel/src/cmd/coel
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

[Here](https://coel-lang.gitbooks.io/coel-programming-language/).

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
