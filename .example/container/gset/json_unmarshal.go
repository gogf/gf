package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/container/gset"
)

func main() {
	b := []byte(`{"Id":1,"Name":"john","Scores":[100,99,98]}`)
	type Student struct {
		Id     int
		Name   string
		Scores *gset.IntSet
	}
	s := Student{}
	json.UnmarshalUseNumber(b, &s)
	fmt.Println(s)
}
