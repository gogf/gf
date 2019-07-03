package main

import (
<<<<<<< HEAD
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
=======
	"github.com/gogf/gf/g/os/glog"

	"github.com/gogf/gf/g/os/gcache"
)

func localCache() {
	result := gcache.GetOrSetFunc("test.key.1", func() interface{} {
		return nil
	}, 1000*60*2)
	if result == nil {
		glog.Error("未获取到值")
	} else {
		glog.Infofln("result is $v", result)
	}
}

func TestCache() {
	for i := 0; i < 100; i++ {
		localCache()
	}
}

func main() {
	TestCache()
>>>>>>> c90ed0d4242527435a3b4c9d7c27742d29c9aaa1
}
