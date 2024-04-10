# RateStor: A Go Rate Limiting Library

RateStor is a simple, efficient, and thread-safe rate-limiting library for Go. It allows you to limit the rate of requests.


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

    for i:=0; i < 101; i++ {
        err := stor.Allow("key", 1*time.Minute, 100)
        if err == ratestor.ErrRateLimitExceeded {
            fmt.Println("Rate limit is exceded")
        }
    }
}
```

## Contributing

Contributions to Wasabi are welcome! Please submit a pull request or create an issue to contribute.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
