package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Println(strings.Trim(`  1  `, "./- \n\r"))
    //fmt.Println(math.MaxInt64)
    //fmt.Println(gtime.Second())
    //fmt.Println(gtime.Nanosecond())
}