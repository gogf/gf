package main

import (
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gogf/gf/g"
)

func main() {
	r, err := g.DB().GetAll(`SELECT TOP 10  * FROM KF_PatInfo_Emergency`)
	fmt.Println(err)
	g.Dump(r.ToList())
}
