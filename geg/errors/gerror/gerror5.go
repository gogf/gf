package main

import (
	"errors"

	"github.com/gogf/gf/g/os/glog"

	"github.com/gogf/gf/g/errors/gerror"
)

func Error1() error {
	return errors.New("test1")
}

func Error2() error {
	return gerror.New("test2")
}

func main() {
	glog.Println(Error1())
	glog.Println(Error2())
}
