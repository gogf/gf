package main

import (
	"encoding/json"
	"fmt"
	"github.com/jin502437344/gf/container/gtype"
)

func main() {
	type Student struct {
		Id     *gtype.Int
		Name   *gtype.String
		Scores *gtype.Interface
	}
	s := Student{
		Id:     gtype.NewInt(1),
		Name:   gtype.NewString("john"),
		Scores: gtype.NewInterface([]int{100, 99, 98}),
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
}
