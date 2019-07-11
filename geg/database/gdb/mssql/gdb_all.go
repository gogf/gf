package main

import (
	"fmt"

	//_ "github.com/denisenkom/go-mssqldb"
	"github.com/gogf/gf/g"
)

func main() {
	r, err := g.DB().GetAll(`SELECT * FROM (SELECT TOP 10 * FROM (SELECT TOP 10  * FROM KF_PatInfo_Emergency  WHERE Report_BZ = 1 AND Examine_BZ = 0 ) as TMP1_ ) as TMP2_`)
	fmt.Println(err)
	g.Dump(r.ToList())
}
