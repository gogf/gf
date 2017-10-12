package main

import (
    "g/os/gfile"
    "fmt"
)



func main() {
    fmt.Println(gfile.PutContents("/tmp/123/1/1/1/1/1/test", "12345678"))
}