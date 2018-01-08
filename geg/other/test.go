package main

import (
    "fmt"
    "reflect"
)


type T struct {
    name string
}

func (t *T)Test() {
    fmt.Println(t.name)
}

func main() {
    t := &T{"john"}
    //fmt.Printf("%p\n", t.Test)
    //fmt.Printf("%p\n", reflect.ValueOf(t).MethodByName("Test").Interface().(func()))
    reflect.ValueOf(t).MethodByName("Test").Interface().(func())()
}