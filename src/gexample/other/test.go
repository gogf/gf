package main

import (
    "fmt"
    "g/os/gfile"
)


func main() {

    path := "/tmp/192.168.2.102.graft.db2"
    fmt.Println(gfile.Exists(path))
}