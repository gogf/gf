package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
)

func main() {
	body := "{\"id\": 413231383385427875}"
	if dat, err := gjson.DecodeToJson(body); err == nil {
		fmt.Println(dat.MustToJsonString())
	}
}
