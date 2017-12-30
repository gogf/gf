package main

import (
    "gitee.com/johng/gf/g/net/gsession"
    "fmt"
)


func main() {
    id := 0
    for i := 0; i < 10; i++ {
        s  := gsession.Get("1")
        if r := s.Get("id"); r != nil {
            id = r.(int)
        }
        id++
        s.Set("id", id)

        fmt.Println(id)
    }

}