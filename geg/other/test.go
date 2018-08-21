package main

import "fmt"

type A struct {
    S string
}

type B struct {
    A
}


func (a *A) editA () {
    a.S += "a"
}

func (b *B) editB () {
    b.S += "b"
}

func main() {
    b := new(B)
    b.editA()
    b.editB()
    b.A.editA()
    fmt.Println(b.S)
}