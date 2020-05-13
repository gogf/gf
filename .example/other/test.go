package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"math"
)

func main() {
	body := "{\"id\": 413231383385427875}"
	fmt.Println(math.MaxFloat32)
	if dat, err := gjson.DecodeToJson(body); err == nil {
		fmt.Println(dat.MustToJsonString())
	}
}
