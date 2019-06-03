package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/g/encoding/gparser"
)

type User struct {
	Id       int  `json:"id" gconv:"i_d"`
}

func main() {
	user := User{100}
	jsonBytes, _ := json.Marshal(user)
	fmt.Println(string(jsonBytes))

	b, _ := gparser.VarToJson(user)
	fmt.Println(string(b))
}