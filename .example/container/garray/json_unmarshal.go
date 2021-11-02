package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
)

func main() {
	b := []byte(`{"Id":1,"Name":"john","Scores":[100,99,98]}`)
	type Student struct {
		Id     int
		Name   string
		Scores *garray.IntArray
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)
}
