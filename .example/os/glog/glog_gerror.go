package main

import (
	"errors"
	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/errors/gerror"
)

func MakeError() error {
	return errors.New("connection closed with normal error")
}

func MakeGError() error {
	return gerror.New("connection closed with gerror")
}

func TestGError() {
	err1 := MakeError()
	err2 := MakeGError()
	g.Log().Error(err1)
	g.Log().Error(err2)
}

func main() {
	TestGError()
}
