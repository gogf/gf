package main

import (
	"errors"

	"github.com/gogf/gf/v2/os/glog"

	"github.com/gogf/gf/v2/errors/gerror"
)

func Error1() error {
	return errors.New("test1")
}

func Error2() error {
	return gerror.New("test2")
}

func main() {
	glog.Print(Error1())
	glog.Print(Error2())
}
