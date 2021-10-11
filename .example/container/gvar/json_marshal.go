package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	type Student struct {
		Id     *g.Var
		Name   *g.Var
		Scores *g.Var
	}
	s := Student{
		Id:     g.NewVar(1),
		Name:   g.NewVar("john"),
		Scores: g.NewVar([]int{100, 99, 98}),
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
}
