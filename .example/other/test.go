package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
)

func main() {

	type Item struct {
		Title string `json:"title"`
		Key   string `json:"key"`
	}

	type M struct {
		Id    string                 `json:"id"`
		Me    map[string]interface{} `json:"me"`
		Txt   string                 `json:"txt"`
		Items []*Item                `json:"items"`
	}

	txt := `{
  "id":"88888",
  "me":{"name":"mikey","day":"20009"},
  "txt":"hello",
  "items":null
 }`

	json, _ := gjson.LoadContent(txt)
	fmt.Println(json)
	m := new(M)
	e := json.ToStructDeep(m)
	fmt.Println(e)
	fmt.Println(m)

}
