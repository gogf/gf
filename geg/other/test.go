package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    type B struct {
        Name string
    }
    type A struct {
        Name  string
        Child B
    }

    a := A {
        Name  : "A",
        Child : B {
            Name : "B",
        },
    }
    b, _ := json.Marshal(a)
    a2 := new(A)
    json.Unmarshal(b, a2)
    fmt.Println(*a2)
}

