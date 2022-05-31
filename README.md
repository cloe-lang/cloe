# Cloe

[![GitHub Action](https://img.shields.io/github/workflow/status/cloe-lang/cloe/test?style=flat-square)](https://github.com/cloe-lang/cloe/actions)
[![Coveralls](https://img.shields.io/coveralls/cloe-lang/cloe/master.svg?style=flat-square)](https://coveralls.io/github/cloe-lang/cloe)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloe-lang/cloe?style=flat-square)](https://goreportcard.com/report/github.com/cloe-lang/cloe)
[![License](https://img.shields.io/github/license/cloe-lang/cloe.svg?style=flat-square)](https://opensource.org/licenses/MIT)

<div align="center">
  <img src="https://raw.githubusercontent.com/cloe-lang/icon/master/spaced.png" alt="logo"/>
</div>

Cloe is the _timeless_ functional programming language.
It aims to be simple and practical.

## Features

- Functional programming
- Immutable data
- Lazy evaluation
- Implicit parallelism, concurrency, and reactiveness

## Installation

```
go get -u github.com/cloe-lang/cloe/...
```

Go 1.8+ is required.

## Documentation

[Here](https://cloe-lang.org).

## Examples

### Hello, world!

```
(print "Hello, world!")
```

### HTTP server

```
(import "http")

(def (handler request)
  ((@ request "respond") "Hello, world!"))

(let requests (http.getRequests ":8080"))

..(map handler requests)
```

See [examples](examples) directory for more.

## License

[MIT](LICENSE)
