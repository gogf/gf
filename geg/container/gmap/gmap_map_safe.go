package main

import (
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g"
)

func main() {
    m := gmap.New()
    m.Set("1", "1")

    m1 := m.Clone()
    m1["2"] = "2"

    g.Dump(m.Clone())
    g.Dump(m1)
}
