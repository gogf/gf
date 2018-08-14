package main

import (
    "fmt"
    "path/filepath"
    "os"
)



func main() {
    fmt.Println(filepath.Abs(filepath.Dir(os.Args[0])))
}