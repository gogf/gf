package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregex"
)

func main() {
    s := `[0-9][0-9]/Jan/[0-9][0-9][0-9][0-9]:[0-9][0-9]:[0-9][0-9]:[0-9][0-9] \+[0-9][0-9][0-9][0-9]`
    s,_  = gregex.ReplaceString(`[A-Za-z]`, `[A-Za-z]`, s)
    fmt.Println(s)
}