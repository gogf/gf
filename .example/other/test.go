package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
)

type A struct {
	D string
	E string
}
type B struct {
	A `json:"a"`
	F string
}

func SystemJsonEncode(a interface{}) string {
	js, err := json.Marshal(a)
	if err != nil {
		return "{}"
	} else {
		return fmt.Sprintf("%s", js)
	}
}

func main() {
	var b B
	fmt.Println(SystemJsonEncode(b))
	fmt.Println(gjson.New(b).MustToJsonString())
}
