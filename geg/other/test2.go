package main

import "fmt"

type User struct {
    Uid   int
}

func New() *User {
    return &User{
        100,
    }
}

func (user *User) Clear() {
    user = New()
}

func main() {
    user := New()
    user.Uid = 10000
    fmt.Println(user)
    user.Clear()
    fmt.Println(user)
}