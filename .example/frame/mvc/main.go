package main

import (
	"fmt"
	"github.com/gogf/gf/.example/frame/mvc/app/model/defaults"
	"github.com/gogf/gf/database/gdb"
)

func main() {
	u := defaults.User{Id: 1, Nickname: "test"}
	fmt.Println(gdb.GetWhereConditionOfStruct(&u))
	fmt.Println(u.Replace())
}
