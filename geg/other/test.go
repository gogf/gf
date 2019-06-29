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

	//fmt.Printf("%+v", err)
	//fmt.Println("==============")
	fmt.Println(gerror.Stack(err))
	return

	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			fmt.Println("ERROR:", err)
			if err, ok := err.(stackTracer); ok {
				for _, f := range err.StackTrace() {
					fmt.Printf("%+s:%d\n", f, f)
				}
			}
			break
		}
		fmt.Println("ERROR:", err)
		if err, ok := err.(stackTracer); ok {
			for _, f := range err.StackTrace() {
				fmt.Printf("%+s:%d\n", f, f)
			}
		}
		fmt.Println()
		err = cause.Cause()
	}

	//err, ok := Test2().(stackTracer)
	//if !ok {
	//	panic("oops, err does not implement stackTracer")
	//}
	//
	//st := err.StackTrace()
	//fmt.Printf("%+v\n", st)
	//fmt.Println()
	//fmt.Printf("%+v\n", Test2())
}
