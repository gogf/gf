package main

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g"
)


func main () {
    a := garray.NewIntArray(0, 0)
    a.Append(1, 2, 3)

    v   := a.Slice()
    v[0] = 4

    g.Dump(a.Slice())
    g.Dump(v)
}
