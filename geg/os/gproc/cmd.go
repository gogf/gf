package main

import (
    "os"
    "fmt"
    "time"
    "os/exec"
)

func main () {
    cmd := exec.Command(os.Args[0], "1")
    time.Sleep(3*time.Second)
    fmt.Println(cmd.Start())
    time.Sleep(time.Hour)
}
