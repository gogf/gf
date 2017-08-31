package main

import (
    "fmt"
    "sort"
)



func main() {
    a := []string{"linux","windows", "windows10"}
    //fmt.Println(a)
    //sort.Strings(a)
    //fmt.Println(a)
    fmt.Println(sort.SearchStrings(a, "windows10"))
}