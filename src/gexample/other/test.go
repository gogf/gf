package main

import (
    "g/os/gfile"
    "fmt"
)


func main() {
    //err := gfile.Mkdir("/tmp/a/b/c/d")
    fmt.Println(gfile.Readable("/"))
    fmt.Println(gfile.Writable("/root"))
}