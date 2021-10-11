package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)
	result, err := db.Model("pw_passageway m,pw_template t").Data("t.status", 99).Where("m.templateId=t.id AND m.status = 0").Update()
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}
