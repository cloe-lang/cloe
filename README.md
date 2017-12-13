# Tisp

[![Circle CI](https://img.shields.io/circleci/project/github/tisp-lang/tisp.svg?style=flat-square)](https://circleci.com/gh/tisp-lang/tisp)
[![codeclimate](https://img.shields.io/codeclimate/github/kabisaict/flow.svg?style=flat-square)](https://codeclimate.com/github/tisp-lang/tisp)
[![Go Report Card](https://goreportcard.com/badge/github.com/tisp-lang/tisp?style=flat-square)](https://goreportcard.com/report/github.com/tisp-lang/tisp)
[![Coveralls](https://img.shields.io/coveralls/tisp-lang/tisp.svg?style=flat-square)](https://coveralls.io/github/tisp-lang/tisp)
[![License](https://img.shields.io/github/license/tisp-lang/tisp.svg?style=flat-square)](https://opensource.org/licenses/MIT)

<div align="center">
  <img src="https://raw.githubusercontent.com/tisp-lang/icon/master/icon.png" alt="logo"/>
</div>

Tisp is a functional programming language.
It aims to be simple, canonical, and practical.

## Installation

```
go get github.com/tisp-lang/tisp/src/cmd/tisp
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

[Here](https://tisp-lang.gitbooks.io/tisp-programming-language/).

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
