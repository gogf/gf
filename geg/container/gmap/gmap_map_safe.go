package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/container/gmap"
)

func main() {
	m := gmap.New()
	m.Set("1", "1")

	m1 := m.Map()
	m1["2"] = "2"

	g.Dump(m.Clone())
	g.Dump(m1)
}
