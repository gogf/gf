package main

import (
	"errors"
	"fmt"

	"github.com/jin502437344/gf/errors/gerror"
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
	fmt.Printf("%s, %-s, %+s\n", err1, err1, err1)
	fmt.Printf("%v, %-v, %+v\n", err1, err1, err1)
	fmt.Printf("%s, %-s, %+s\n", err2, err2, err2)
	fmt.Printf("%v, %-v, %+v\n", err2, err2, err2)
}
