package main

import (
    "fmt"
    "gf"
)

func main() {
    r,_ := gf.Console.Value.GetIndex(1)
    fmt.Printf("%s", r)
}
