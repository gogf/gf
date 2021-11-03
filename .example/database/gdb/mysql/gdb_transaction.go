package main

import (
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	var (
		db    = g.DB()
		table = "user"
	)
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	if err = tx.Begin(); err != nil {
		panic(err)
	}
	_, err = tx.Model(table).Data(g.Map{"id": 1, "name": "john"}).Insert()
	if err = tx.Rollback(); err != nil {
		panic(err)
	}
	_, err = tx.Model(table).Data(g.Map{"id": 2, "name": "smith"}).Insert()
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}
