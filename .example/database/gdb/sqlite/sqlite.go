package main

import (
	"fmt"

	"github.com/jin502437344/gf/database/gdb"
	"github.com/jin502437344/gf/frame/g"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Name: "/tmp/my.db",
				Type: "sqlite",
			},
		},
	})
	db := g.DB()
	if db == nil {
		panic("db create failed")
	}

	// 创建表
	sql := `CREATE TABLE user (
       uid  INT PRIMARY KEY NOT NULL,
       name VARCHAR(30) NOT NULL
   );`
	if _, err := db.Exec(sql); err != nil {
		fmt.Println(err)
	}

	// 写入数据
	result, err := db.Table("user").Data(g.Map{"uid": 1, "name": "john"}).Save()
	if err == nil {
		fmt.Println(result.RowsAffected())
	} else {
		fmt.Println(err)
	}

	// 删除表
	sql = `DROP TABLE user;`
	if _, err := db.Exec(sql); err != nil {
		fmt.Println(err)
	}
}
