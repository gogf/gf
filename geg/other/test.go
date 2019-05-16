package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/util/gconv"
)

func main() {
	fmt.Println(gfile.Dir("/"))
	return
	a := []int{1,2,3}
	fmt.Println(a[:0])
	return
	type Person struct{
		Name string
	}
	type Staff struct{
		Person
		StaffId int
	}
	staff  := &Staff{}
	params := g.Map{
		"Name"    : "john",
		"StaffId" : "10000",
	}
	gconv.Struct(params, staff)
	fmt.Println(staff)
}