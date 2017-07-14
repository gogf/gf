package main

import (
    "fmt"
    "regexp"
)



func main() {
    str := "Welcome for Beijing-Tianjin?CRH train."
    reg := regexp.MustCompile("\\?")
    index := 0
    fmt.Println(reg.ReplaceAllStringFunc(str, func (s string) string {
        index ++
        return fmt.Sprintf("$%d", index)
    }))
}