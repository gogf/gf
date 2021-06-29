package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	var (
		err   error
		db    = g.DB()
		table = "user"
	)
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
		}
	}()
	if _, err = tx.Model(table).Data(g.Map{"id": 1, "name": "john"}).Insert(); err != nil {
		panic(err)
	}
	if err = tx.SavePoint("MyPoint"); err != nil {
		panic(err)
	}
	if _, err = tx.Model(table).Data(g.Map{"id": 2, "name": "smith"}).Insert(); err != nil {
		panic(err)
	}
	if _, err = tx.Model(table).Data(g.Map{"id": 3, "name": "green"}).Insert(); err != nil {
		panic(err)
	}
	if err = tx.RollbackTo("MyPoint"); err != nil {
		panic(err)
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}
