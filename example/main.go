package main

import (
	stderrors "errors"
	"fmt"
	"github.com/mawngo/go-errors"
)

var ErrUhOh = stderrors.New("uh oh")

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
