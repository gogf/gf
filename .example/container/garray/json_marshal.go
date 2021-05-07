package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/container/garray"
)

func main() {
	type Student struct {
		Id     int
		Name   string
		Scores *garray.IntArray
	}
	s := Student{
		Id:     1,
		Name:   "john",
		Scores: garray.NewIntArrayFrom([]int{100, 99, 98}),
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
}
