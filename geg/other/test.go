package main

import (
    "fmt"
    "reflect"
)

type Test struct {
    name string
}

func (t *Test) Show(s string) {
    fmt.Println(t.name)
    fmt.Println(s)
}

type F func(string)

func main() {
    t := &Test{}
    t.name = "john"
    reflect.ValueOf(t).Method(0).Interface().(F)("123")

}