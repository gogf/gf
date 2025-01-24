package tpl

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

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

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
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
	outputDir := "./output"
	tplRootDir := "./testdata"
	tplRootDir = gfile.Abs(tplRootDir)
	fmt.Println(tplRootDir)
	tplList, err := gfile.ScanDirFile(tplRootDir, "*.tpl", true)
	if err != nil {
		panic(err)
	}
	fmt.Println(tplList)

	fmt.Printf("%#v\n", Table{})
	fmt.Printf("%#v\n", TableField{})

	tables := GetTables(ctx, db)
	view := gview.New()

	for _, table := range tables {
		table.PackageName = "github.com/gogf/gf/cmd/gf/v2"
		tplData := g.Map{
			"table":  table,
			"tables": tables,
		}
		fmt.Println(table.FieldsJsonStr("Snake"))
		for _, tpl := range tplList {
			tplDir := gfile.Dir(tpl)

			res, err := view.Parse(ctx, tpl, tplData)
			if err != nil {
				panic(err)
			}
			// fmt.Println(res, err)
			filePath := filepath.FromSlash(fmt.Sprintf(outputDir+"%s/%s.go", strings.TrimPrefix(tplDir, tplRootDir), table.Name))
			err = gfile.PutContents(filePath, res)
			if err != nil {
				panic(err)
			}
			utils.GoFmt(filePath)
		}

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
