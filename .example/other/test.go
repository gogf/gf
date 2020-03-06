package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type MyInt int

//func (i *MyInt) UnmarshalValue(interface{}) error {
//	*i = 10
//	return nil
//}
func main() {
	type User struct {
		Id MyInt
	}
	user := new(User)
	err := gconv.Struct(g.Map{
		"id": 1,
	}, user)
	fmt.Println(err)
	fmt.Println(user)
}
