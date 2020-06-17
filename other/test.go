package main

import (
	"fmt"
)

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	body := "{\"id\": 413231383385427875}"
	m := make(map[string]interface{})
	json.Unmarshal([]byte(body), &m)
	b, _ := json.Marshal(m)
	fmt.Println(string(b))
}
