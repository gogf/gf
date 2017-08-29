package main

import (
    "os/exec"
    "fmt"
)



func main() {
    b , err :=exec.Command("sh", "-c", "ls /home").Output()
    fmt.Println(b)
    fmt.Println(err)
}