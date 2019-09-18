package main

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/text/gstr"
)

func main() {
	m, _ := gstr.Parse("map[a]=1&map[b]=2")
	g.Dump(m)
	fmt.Println(reflect.TypeOf(m["map"].(map[string]interface{})["b"]))
}
