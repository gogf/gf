package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
)

type A struct {
	F string
	G string
}

type B struct {
	A
	H string
}

type C struct {
	A A
	F string
}

type D struct {
	I A
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
	fmt.Println("encoding/json", SystemJsonEncode(B{}))
	fmt.Println("gjson", gjson.New(B{}).MustToJsonString())
	fmt.Println()
	fmt.Println("encoding/json", SystemJsonEncode(C{}))
	fmt.Println("gjson", gjson.New(C{}).MustToJsonString())
	fmt.Println()
	fmt.Println("encoding/json", SystemJsonEncode(D{}))
	fmt.Println("gjson", gjson.New(D{}).MustToJsonString())
}
