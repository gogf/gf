package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gconv"
)

func main() {
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