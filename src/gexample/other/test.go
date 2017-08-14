package main

import (
    "fmt"
    "g/encoding/gmd5"
    "errors"
)


func main() {
    fmt.Println(gmd5.Encode(1))
    fmt.Println(gmd5.Encode(errors.New("123")))

}