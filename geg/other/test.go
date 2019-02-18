package main

import (
    "fmt"
)

func main() {
    s := "abc我是中国人é"
    fmt.Println(len(s))

    for i := 0; i < len(s); i++ {
        fmt.Println(s[i])
    }
}