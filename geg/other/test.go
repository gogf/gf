package main

import "fmt"

const (
    CREATE = 1 << iota
    WRITE
    REMOVE
    RENAME
    CHMOD
)


func main(){

    fmt.Println(CREATE)
    fmt.Println(WRITE)
    fmt.Println(REMOVE)
    fmt.Println(RENAME)
}