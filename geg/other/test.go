package main
import (
    "reflect"
    "fmt"
)

type B struct {
    Name string
}

func(b *B) Test() {

}

func main() {
    b := &B{}
    t := reflect.ValueOf(b).Elem().Type()
    //n := reflect.New(t)
    fmt.Println(t.NumMethod())
}