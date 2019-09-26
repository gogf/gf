package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	//type User struct {
	//	Uid  int
	//	Name *gvar.Var
	//}

	//user := new(User)
	////user.Name = g.NewVar("john")
	//g.Dump(gconv.Map(user))

	_, e := db.Table("test").Data(g.Map{
		"name": nil,
	}).Update()
	if e != nil {
		panic(e)
	}

}
