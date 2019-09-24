package main

import (
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/container/gmap"
)

func main() {
	m := gmap.NewIntIntMap()
	m.Set(1, 2)
	m.Set(3, 4)
	b, err := json.Marshal(m)
	fmt.Println(err)
	fmt.Println(string(b))
}
