package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	type Student struct {
		Id     int
		Name   string
		Scores *glist.List
	}
	s := Student{
		Id:     1,
		Name:   "john",
		Scores: glist.NewFrom(g.Slice{100, 99, 98}),
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
}
