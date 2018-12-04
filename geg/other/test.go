package main

import (
	"fmt"
	"gitee.com/johng/gf/g/encoding/gparser"
)

func main() {
	type User struct {
		Uid  int    `json:"uid"`
		Name string `json:"name"`
	}
	user := User{1, "john"}
	b, err := gparser.VarToJson(user)
	fmt.Println(err)
	fmt.Println(string(b))
}
