package main

import (
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/g/container/gmap"
)

func main() {
	//m := gmap.New()
	//m.Set("k", "v")
	//m.Set("1", "2")
	//b, err := json.Marshal(m)
	//fmt.Println(err)
	//fmt.Println(string(b))

	m := gmap.NewIntIntMap()
	m.Set(1, 2)
	m.Set(3, 4)
	b, err := json.Marshal(m)
	fmt.Println(err)
	fmt.Println(string(b))
}
