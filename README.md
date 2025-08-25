# Go Errors

Simple errors with stack trace support.

Copied from [Thanos](https://github.com/thanos-io/thanos/tree/main/pkg/errors) with some modification to meet my need.

```
test_wrapper
> go-errors.TestWrapf	E:/Dev/Golang/go-errors/errors_test.go:41
> testing.tRunner	C:/Program Files/Go/src/testing/testing.go:1792
> runtime.goexit	C:/Program Files/Go/src/runtime/asm_amd64.s:1700
test_error_message
> go-errors.TestWrapf	E:/Dev/Golang/go-errors/errors_test.go:40
> testing.tRunner	C:/Program Files/Go/src/testing/testing.go:1792
> runtime.goexit	C:/Program Files/Go/src/runtime/asm_amd64.s:1700
```

## Installation

Require go 1.25+

```shell
go get -u github.com/mawngo/go-errors
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/mawngo/go-errors"
)

var ErrUhOh = errors.New("uh oh")

func main() {
	wrap := errors.Wrapf(ErrUhOh, "uhoh wrapped")
	fmt.Printf("%+v", wrap)
	//nolint
	//uhoh wrapped
	//> main.main	E:/Dev/Golang/go-errors/example/main.go:11
	//> runtime.main	C:/Program Files/Go/src/runtime/proc.go:283
	//> runtime.goexit	C:/Program Files/Go/src/runtime/asm_amd64.s:1700
	//uh oh
}
```