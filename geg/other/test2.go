package main

import (
    "gitee.com/johng/gf/g/encoding/gjson"
)

func main() {
	j := gjson.New(nil)
	j.Set("array", []int{1,2,3})
    j.Append("array", 4)
    j.Dump()
}