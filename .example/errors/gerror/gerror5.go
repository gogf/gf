package main

import (
	"errors"

	"github.com/jin502437344/gf/os/glog"

	"github.com/jin502437344/gf/errors/gerror"
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
