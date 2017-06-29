package main

import "fmt"
import "reflect"

type st struct{
    age  int
    name string
}


func (_ st)Echo(str string){
    fmt.Printf("echo(%s)\n", str)
}

func main() {
    s := st {age: 18, name:"john"}
    p := reflect.ValueOf("hallo")
    v := reflect.ValueOf(s)
    v.MethodByName("Echo").Call([]reflect.Value{ p })
    fmt.Println(v.FieldByName("name"))
}