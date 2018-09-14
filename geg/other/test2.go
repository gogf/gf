package main

import (
    "gitee.com/johng/gf/g/encoding/gjson"
)

func main() {
	j := gjson.NewUnsafe([]int{1,2,3})
    j.Append("", 4)
    j.Append("", "abc")
    j.Dump()

}