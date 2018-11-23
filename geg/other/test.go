package main

import "fmt"

func main() {
    s := "我是中国人//"
    for _, v := range s {
        fmt.Println(v)
    }
    fmt.Println(s)
}