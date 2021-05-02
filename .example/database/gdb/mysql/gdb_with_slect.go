package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gmeta"
)

func main() {
	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScore struct {
		gmeta.Meta `orm:"table:user_score"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user"`
		Id         int          `json:"id"`
		Name       string       `json:"name"`
		UserDetail *UserDetail  `orm:"with:uid=id"`
		UserScores []*UserScore `orm:"with:uid=id"`
	}

	db := g.DB()
	var user *User
	err := db.Model(user).WithAll().Where("id", 3).Scan(&user)
	if err != nil {
		panic(err)
	}
	g.Dump(user)
}
