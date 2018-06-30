package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gjson"
)

func main() {
    content := `[0.00000000059, 1.598877777409]`
    j, _ := gjson.LoadContent([]byte(content), "json")
    fmt.Println(j.GetString("0"))
}
