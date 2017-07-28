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

    fmt.Println(5.000/2)

}