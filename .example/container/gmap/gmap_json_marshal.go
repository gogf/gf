package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/container/gmap"
)

func main() {
	m := gmap.New()
	m.Sets(g.MapAnyAny{
		"name":  "john",
		"score": 100,
	})
	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}
