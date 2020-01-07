package main

import "fmt"

func main() {
	type User struct {
		Id int
	}
	u1 := &User{1}
	u2 := *u1
	u2.Id = 2
	fmt.Println(u1)
	fmt.Println(u2)
}
