package main

import (
	"errors"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
)

func Error1() error {
	return errors.New("test1")
}

func Error2() error {
	return gerror.New("test2")
}

func main() {
	err1 := Error1()
	err2 := Error2()
	fmt.Println("err1:\n", gerror.Stack(err1))
	fmt.Println("err2:\n", gerror.Stack(err2))
}
