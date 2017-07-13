package main

import (
    "fmt"
    "g/os/gconsole"
)

func doEcho() {
    fmt.Println("do echo")
}

func main() {
    fmt.Println(gconsole.Value.GetAll())

    fmt.Println(gconsole.Value.GetIndex(1))

    gconsole.BindHandle("echo", doEcho)
    gconsole.RunHandle("echo")

    gconsole.AutoRun()
}
