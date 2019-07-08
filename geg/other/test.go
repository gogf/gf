package main

import (
	"fmt"

	"github.com/gogf/gf/g"
)

type User struct {
	Uid  int
	Name string
}

func main() {
	if r, err := g.DB().Table("user").Where("uid=?", 1).One(); r != nil {
		u := new(User)
		if err := r.ToStruct(u); err == nil {
			fmt.Println(" uid:", u.Uid)
			fmt.Println("name:", u.Name)
		} else {
			fmt.Println(err)
		}
	} else if err != nil {
		fmt.Println(err)
	}
}
