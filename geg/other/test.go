package main

import (
    "fmt"
    "reflect"
    "runtime"
)

type T struct {

}
func (t *T) Test() {

}

func main() {
    t := new(T)
    fmt.Println(runtime.FuncForPC(reflect.ValueOf(t.Test).Pointer()).Name())
}
