package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	b := []byte(`{"Id":1,"Name":"john","Scores":[100,99,98]}`)
	type Student struct {
		Id     *g.Var
		Name   *g.Var
		Scores *g.Var
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)
}
