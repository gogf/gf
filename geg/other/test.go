package main

import (
    "fmt"
    "reflect"
)

type T struct {

}

func (t *T) Test2Test() {}

func main() {
    obj := &T{}
    v := reflect.ValueOf(obj).MethodByName("Test2Test")
    fmt.Println(v.IsValid())
}