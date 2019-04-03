package main

import (
	"fmt"
	"github.com/gogf/gf/g/encoding/gjson"
)

func main() {
	config := `
v1    = 1
v2    = "true"
v3    = "off"
v4    = "1.23"
array = [1,2,3]
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
	j, err := gjson.LoadContent(config)
	fmt.Println(err)
	fmt.Println(j.ToJsonIndentString())
}
