package main

import (
    "fmt"
)


type T struct {
    name string
}


func (t *T) swap(t2 *T) {
    *t = &t2
}

func main() {
    t1 := &T{"john"}
    t2 := &T{"smith"}
    t2.swap(t2)

    fmt.Println(t1)
}