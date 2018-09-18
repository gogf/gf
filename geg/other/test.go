package main

import (
    "gitee.com/johng/gf/g/container/gmap"
    "fmt"
)

func main() {
	m1 := gmap.NewIntInterfaceMap()
	m1.Set(1, gmap.NewIntIntMap())
	v1 := m1.Get(1).(*gmap.IntIntMap)
	v1.Set(2, 2)
	fmt.Println(m1.Get(1).(*gmap.IntIntMap).Size())
}