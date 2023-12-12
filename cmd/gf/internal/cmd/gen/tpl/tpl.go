package tpl

import (
	"context"
	"fmt"

	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
)

// init description
//
// createTime: 2023-10-23 15:42:14
//
// author: hailaz
func init() {
	nodeDefault := gdb.ConfigNode{
		Link: fmt.Sprintf("mysql:root:%s@tcp(127.0.0.1:3306)/focus?loc=Local&parseTime=true", "root123"),
	}

	gdb.AddConfigNode("test", nodeDefault)
}

// GetTables description
//
// createTime: 2023-12-11 16:21:43
//
// author: hailaz
func GetTables(ctx context.Context, db gdb.DB) Tables {
	tablesName, err := db.Tables(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(tablesName)
	tables := make(Tables, 0)
	for _, v := range tablesName {
		t, err := NewTable(ctx, db, v)
		if err != nil {
			panic(err)
		}
		t.SortFields(true)
		tables = append(tables, t)
	}
	return tables
}

func Tpl() {
	ctx := context.TODO()
	db, _ := gdb.Instance("test")

	fmt.Printf("%#v\n", Table{})
	fmt.Printf("%#v\n", TableField{})

	tables := GetTables(ctx, db)
	view := gview.New()
	for _, table := range tables {
		res, err := view.Parse(ctx, "./testdata/dao.tpl",
			g.Map{
				"table":  table,
				"tables": tables,
			},
		)
		if err != nil {
			panic(err)
		}
		// fmt.Println(res, err)
		gfile.PutContents(fmt.Sprintf("./output/%s.go", table.Name), res)
	}

	return
	// 	res, err := view.ParseContent(ctx, `{{range $i,$v := .tables}}{{$i}} ---- {{$v}} {{$v.Name}}
	// {{$v.FieldsJsonStr}}
	// 	{{range $k,$vv := $v.Fields}}{{$k}} ---- {{$vv}}
	// 	{{end}}
	// {{end}}

	// `, g.Map{"tables": tables})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(res, err)

}
