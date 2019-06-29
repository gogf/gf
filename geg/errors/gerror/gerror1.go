package main

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/gogf/gf/g/errors/gerror"
)

func Test1() error {
	return gerror.New("test")
}

func Test2() error {
	return gerror.Wrap(Test1(), "error test1")
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func main() {
	err := Test2()
	fmt.Printf("%s\n", err)
	fmt.Printf("%v\n", err)
	fmt.Printf("%+v\n", err)
	fmt.Println(gerror.Stack(err))
	return

}
