package main

import (
    "fmt"
    "reflect"
)

type S struct {

}

func main() {
    s := S{}
    fmt.Println(reflect.ValueOf(s).Kind())
    fmt.Println(reflect.ValueOf(&s).Elem().Kind())

    v := reflect.ValueOf(s).Interface()
    fmt.Println(reflect.ValueOf(v).Kind())
    fmt.Println(reflect.ValueOf(&v).Elem().Type().PkgPath())
}

