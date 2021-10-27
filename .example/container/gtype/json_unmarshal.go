package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/container/gtype"
)

func main() {
	b := []byte(`{"Id":1,"Name":"john","Scores":[100,99,98]}`)
	type Student struct {
		Id     *gtype.Int
		Name   *gtype.String
		Scores *gtype.Interface
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)
}
