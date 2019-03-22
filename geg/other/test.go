package main

import (
    "encoding/json"
    "fmt"
)

func main() {

    fmt.Println(json.Valid([]byte("111")))
}