package main

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
)

func main() {
	a := garray.NewIntArray()
	a.Append(1, 2, 3)

	v := a.Slice()
	v[0] = 4

	g.Dump(a.Slice())
	g.Dump(v)
}
