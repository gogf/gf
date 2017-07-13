package main

import "fmt"


type Student struct {
    Human
    school string
}
type Employer struct {
    Human
    company string
}
type Human struct {
    name  string
    age   int
    phone string
}

//implement Human method
func (h *Human) Show() {
    h.show()
}

func (h *Human) show() {
    fmt.Println("human show")
}

func (s *Student) show() {
    fmt.Println("Student show")
}

//在go中也有方法的重写和继承

func main() {
    s := Student{}
    s.Show()
}