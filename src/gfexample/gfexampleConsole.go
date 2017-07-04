package main

import (
    "fmt"
    "gf"
)

func doEcho() {
    fmt.Println("do echo")
}

func main() {
    fmt.Println(gf.Console.Value.GetAll())

    fmt.Println(gf.Console.Value.GetIndex(1))

    gf.Console.BindHandle("echo", doEcho)
    gf.Console.RunHandle("echo")

    gf.Console.AutoRun()
}
