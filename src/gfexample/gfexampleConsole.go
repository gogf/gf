package main

import (
    "fmt"
    "gf"
)

func doEcho() {
    fmt.Println("do echo")
}

func main() {
    r,_ := gf.Console.Value.GetIndex(1)
    fmt.Printf("%s", r)

    gf.Console.BindHandle("echo", doEcho)
    gf.Console.RunHandle("echo")
}
