package main

import (
	"fmt"
	"github.com/gogf/gf/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)
	result, err := db.Table("pw_passageway m,pw_template t").Data("t.status", 99).Where("m.templateId=t.id AND m.status = 0").Update()
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}
