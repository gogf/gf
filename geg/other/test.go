package main

import (
	"github.com/gogf/gf/g/encoding/gjson"
)

func main() {
	j := gjson.New(`[1,2,3]`)
	j.Remove("1")
	j.Dump()
}
