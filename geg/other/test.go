package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gconv"
	"log"
	"os"
)

func main() {
	var mylog = log.New(os.Stdout, "[Api] ", log.LstdFlags|log.Lshortfile)
	mylog.Println(123)
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