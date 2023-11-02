package tpl

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
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

var (
	defaultTypeMapping = map[DBFieldTypeName]CustomAttributeType{
		"decimal": {
			Type: "float64",
		},
		"money": {
			Type: "float64",
		},
		"numeric": {
			Type: "float64",
		},
		"smallmoney": {
			Type: "float64",
		},
	}
)

type (
	DBFieldTypeName     = string
	CustomAttributeType struct {
		Type   string `brief:"custom attribute type name"`
		Import string `brief:"custom import for this type"`
	}
)

// TableField description
type TableField struct {
	gdb.TableField
	LocalType string
	JsonName  string
}

type TableFields []*TableField

func (s TableFields) Len() int      { return len(s) }
func (s TableFields) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TableFields) Less(i, j int) bool {
	return strings.Compare(s[i].Name, s[j].Name) < 0
}

// Tables description
type Tables struct {
	Name         string
	db           gdb.DB
	Fields       TableFields
	FieldsSource map[string]*gdb.TableField
	Imports      map[string]struct{}
}

// Name description
//
// createTime: 2023-10-23 16:17:30
//
// author: hailaz
func (t *Tables) Show() string {
	return fmt.Sprintf("table name is %s", t.Name)
}

// toTableFields description
//
// createTime: 2023-10-23 17:22:40
//
// author: hailaz
func (t *Tables) toTableFields() {
	t.Fields = make(TableFields, len(t.FieldsSource))
	for _, v := range t.FieldsSource {
		field := &TableField{
			TableField: *v,
		}
		appendImport := field.GetLocalTypeName(context.Background(), t.db, Input{
			TypeMapping: defaultTypeMapping,
			StdTime:     false,
		})
		if appendImport != "" {
			t.Imports[appendImport] = struct{}{}
		}
		t.Fields[v.Index] = field
	}
}

// SortFields description
//
// createTime: 2023-10-23 17:18:22
//
// author: hailaz
func (t *Tables) SortFields(isReverse bool) {
	if isReverse {
		sort.Sort(sort.Reverse(t.Fields))
	} else {
		sort.Sort(t.Fields)
	}
}

// JsonStr description
//
// createTime: 2023-10-23 17:29:39
//
// author: hailaz
func (t *Tables) FieldsJsonStr() string {
	mapStr := make(map[string]interface{}, len(t.Fields))
	for _, v := range t.Fields {
		mapStr[v.Name] = ""
	}
	b, _ := json.Marshal(mapStr)
	return string(b)
}

// Input description
type Input struct {
	StdTime      bool
	GJsonSupport bool
	TypeMapping  map[DBFieldTypeName]CustomAttributeType
}

// GetLocalTypeName description
//
// createTime: 2023-10-25 15:43:06
//
// author: hailaz
func (field *TableField) GetLocalTypeName(ctx context.Context, db gdb.DB, in Input) (appendImport string) {
	var (
		err              error
		localTypeName    gdb.LocalType
		localTypeNameStr string
	)
	if in.TypeMapping != nil && len(in.TypeMapping) > 0 {
		var (
			tryTypeName string
		)
		tryTypeMatch, _ := gregex.MatchString(`(.+?)\((.+)\)`, field.Type)
		if len(tryTypeMatch) == 3 {
			tryTypeName = gstr.Trim(tryTypeMatch[1])
		} else {
			tryTypeName = gstr.Split(field.Type, " ")[0]
		}
		if tryTypeName != "" {
			if typeMapping, ok := in.TypeMapping[strings.ToLower(tryTypeName)]; ok {
				localTypeNameStr = typeMapping.Type
				appendImport = typeMapping.Import
			}
		}
	}

	if localTypeNameStr == "" {
		localTypeName, err = db.CheckLocalTypeForField(ctx, field.Type, nil)
		if err != nil {
			panic(err)
		}
		localTypeNameStr = string(localTypeName)
		switch localTypeName {
		case gdb.LocalTypeDate, gdb.LocalTypeDatetime:
			if in.StdTime {
				localTypeNameStr = "time.Time"
			} else {
				localTypeNameStr = "*gtime.Time"
				appendImport = "github.com/gogf/gf/v2/os/gtime"
			}

		case gdb.LocalTypeInt64Bytes:
			localTypeNameStr = "int64"

		case gdb.LocalTypeUint64Bytes:
			localTypeNameStr = "uint64"

		// Special type handle.
		case gdb.LocalTypeJson, gdb.LocalTypeJsonb:
			if in.GJsonSupport {
				localTypeNameStr = "*gjson.Json"
				appendImport = "github.com/gogf/gf/v2/encoding/gjson"
			} else {
				localTypeNameStr = "string"
			}
		}
	}
	field.LocalType = localTypeNameStr
	field.JsonName = gstr.CaseConvert(field.Name, gstr.Camel)

	return
}

func Tpl() {
	ctx := context.TODO()
	db, _ := gdb.Instance("test")
	tablesName, err := db.Tables(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(tablesName)
	tables := make([]Tables, 0)
	for _, v := range tablesName {
		f, _ := db.TableFields(ctx, v)
		t := Tables{
			Name:         v,
			FieldsSource: f,
			db:           db,
			Imports:      make(map[string]struct{}),
		}
		t.toTableFields()
		t.SortFields(true)
		tables = append(tables, t)
	}

	view := gview.New()
	res, err := view.ParseContent(ctx, `{{range $i,$v := .tables}}{{$i}} ---- {{$v}} {{$v.Name}} 
{{$v.FieldsJsonStr}}
	{{range $k,$vv := $v.Fields}}{{$k}} ---- {{$vv}} 
	{{end}}
{{end}}

`, g.Map{"tables": tables})
	if err != nil {
		panic(err)
	}

	fmt.Println(res, err)

}
