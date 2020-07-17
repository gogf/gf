package main

import (
	"github.com/jin502437344/gf/frame/g"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	tables, e := db.Tables()
	if e != nil {
		panic(e)
	}
	if tables != nil {
		g.Dump(tables)
		for _, table := range tables {
			fields, err := db.TableFields(table)
			if err != nil {
				panic(err)
			}
			g.Dump(fields)
		}
	}
}
