package main

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
)

func main() {
	m := gmap.New()
	m.Set("1", "1")

	m1 := m.Map()
	m1["2"] = "2"

	g.Dump(m.Clone())
	g.Dump(m1)
}
