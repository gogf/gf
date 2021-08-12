package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Log().Line().Debug("this is the short file name with its line number")
	g.Log().Line(true).Debug("lone file name with line number")
}
