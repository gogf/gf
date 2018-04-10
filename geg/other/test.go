package main

import (
    "fmt"
    "github.com/fatih/structs"
)

type T struct {
    Name string
    Age  int
    S    struct {
        N string
    }
}

func (t *T) Test2Test() {}

func main() {
    v := &T{
        Name : "john",
        Age  : 18,
    }
    v.S.N = "ttt"
    fmt.Println(structs.Map(v))

}