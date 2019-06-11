package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/container/garray"
)

func main() {
	a := garray.NewIntArray()
	a.Append(1, 2, 3)

	v := a.Slice()
	v[0] = 4

	g.Dump(a.Slice())
	g.Dump(v)
}
