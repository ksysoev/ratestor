# RateStor: A Go Rate Limiting Library

[![RateStor](https://github.com/ksysoev/ratestor/actions/workflows/main.yml/badge.svg)](https://github.com/ksysoev/ratestor/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/ksysoev/ratestor/graph/badge.svg?token=0TWEWEJW3B)](https://codecov.io/gh/ksysoev/ratestor)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

RateStor is a simple, efficient, and thread-safe rate-limiting library for Go. It allows you to limit the rate of requests.
The library provides a flexible way to define different limits for different keys.

## Installation

To install `ratestor`, use the `go get` command:

```sh
go get github.com/ksysoev/ratestor
```

## Usage 

```go
package main 

import (
    "fmt"
    "time"

    "github.com/ksysoev/ratestor"
)

func main () {
    stor := ratestor.NewRateStor()
    defer stor.Close()

    err := stor.Allow("user1_min", 1*time.Minute, 100)
    if err == ratestor.ErrRateLimitExceeded {
        fmt.Println("Minute limit is exceded")
    }

    err := stor.Allow("user1_10min", 10*time.Minute, 300)
    if err == ratestor.ErrRateLimitExceeded {
        fmt.Println("10 Minute limit is exceded")
    }
}
```

## Contributing

Contributions to Wasabi are welcome! Please submit a pull request or create an issue to contribute.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
