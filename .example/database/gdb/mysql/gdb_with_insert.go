package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gmeta"
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
	err := db.Transaction(gctx.New(), func(ctx context.Context, tx *gdb.TX) error {
		for i := 1; i <= 5; i++ {
			// User.
			user := User{
				Name: fmt.Sprintf(`name_%d`, i),
			}
			lastInsertId, err := db.Ctx(ctx).Model(user).Data(user).OmitEmpty().InsertAndGetId()
			if err != nil {
				return err
			}
			// Detail.
			userDetail := UserDetail{
				Uid:     int(lastInsertId),
				Address: fmt.Sprintf(`address_%d`, lastInsertId),
			}
			_, err = db.Ctx(ctx).Model(userDetail).Data(userDetail).OmitEmpty().Insert()
			if err != nil {
				return err
			}
			// Scores.
			for j := 1; j <= 5; j++ {
				userScore := UserScore{
					Uid:   int(lastInsertId),
					Score: j,
				}
				_, err = db.Ctx(ctx).Model(userScore).Data(userScore).OmitEmpty().Insert()
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	fmt.Println(err)
}
