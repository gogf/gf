package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregex"
    "strings"
)



func main() {
    newWhere := "?????"
    counter  := 0
    newWhere, _ = gregex.ReplaceStringFunc(`\?`, newWhere, func(s string) string {
        counter++
        if counter == 4 {
            return "?" + strings.Repeat(",!", 5 - 1)
        }
        return s
    })
    fmt.Println(newWhere)
}