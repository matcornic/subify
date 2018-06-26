Fluent
=========

Fluent HTTP client for Golang. With timeout, retries and exponential back-off support.

Usage:

```go
package main

import (
  "fmt"
  "github.com/lafikl/fluent"
  "time"
)

func main() {
  req := fluent.New()
  req.Post("http://example.com").
    InitialInterval(time.Duration(time.Millisecond)).
    Json([]int{1, 3, 4}).
    Retry(3)

  res, err := req.Send()

  if err != nil {
    fmt.Println(err)
  }

  fmt.Println("donne ", res)

  // They can be separated if you don't like chaining ;)
  // for example:
  // req.Get("http://example.com")
  // req.Retry(3)
}

```

http://godoc.org/github.com/lafikl/fluent


