package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
)

func main() {
	db := g.DB()
	// 开启调试模式，以便于记录所有执行的SQL
	db.SetDebug(true)

	type User struct {
		Uid  int
		Name string
	}
	user := (*User)(nil)
	fmt.Println(user)
	err := db.Table("test").Where("id=1").Struct(&user)
	fmt.Println(err)
	fmt.Println(user)
}
