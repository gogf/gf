package main

import (
    "fmt"
    //"reflect"
    "reflect"
)
//import "reflect"

type gtInterface interface {
    Run()
}

type st struct {
    age  int
    Name string
}

type mySt struct {
    st
}

func (_ st) Echo(str string) {
    fmt.Printf("echo(%s)\n", str)
}
func (_ *st) Echo2(str string) {
    fmt.Printf("echo2(%s)\n", str)
}

func (_ st) Echo3() {
    fmt.Println("echo3()")
}

func Echo3() {
    fmt.Println("echo3()")
}

type DefaultFunc func()

func Call(i DefaultFunc) {
    i()
    //reflect.ValueOf(i).Call([]reflect.Value{})
}
func main() {
    s  := st {16,"john"}


    //p  := reflect.ValueOf("halloo")
    v  := reflect.ValueOf(s)
    //v2 := reflect.ValueOf(&s)
    //// 调用st结构体的方法
    //v.MethodByName("Echo").Call([]reflect.Value{p})
    //// 我们需要调用的是实体结构体指针的方法，注意v2与v2的区别，以及方法定义的区别
    //v2.MethodByName("Echo2").Call([]reflect.Value{p})
    //v.MethodByName()
    fmt.Println(v.FieldByName("name"))



}