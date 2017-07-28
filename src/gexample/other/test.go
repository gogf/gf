package main

import (
    "fmt"
)

type ttt struct {
    Name string
    Age  int `json:"age"`
    Info struct{
        grade string
    }
}
func main() {

    var a interface{}
    var b struct{}
    fmt.Println(&a)
    fmt.Println(&b)

}