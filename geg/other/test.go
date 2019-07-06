package main

import (
<<<<<<< HEAD
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
=======
	"fmt"
>>>>>>> d0fe2d2f75a91f44d59969bb25fda3729eabcdde

	"github.com/gogf/gf/g"
)

type User struct {
	Uid  int
	Name string
}

func main() {
<<<<<<< HEAD
	TestCache()
>>>>>>> c90ed0d4242527435a3b4c9d7c27742d29c9aaa1
=======
	if r, err := g.DB().Table("user").Where("uid=?", 1).One(); r != nil {
		u := new(User)
		if err := r.ToStruct(u); err == nil {
			fmt.Println(" uid:", u.Uid)
			fmt.Println("name:", u.Name)
		} else {
			fmt.Println(err)
		}
	} else if err != nil {
		fmt.Println(err)
	}
>>>>>>> d0fe2d2f75a91f44d59969bb25fda3729eabcdde
}
