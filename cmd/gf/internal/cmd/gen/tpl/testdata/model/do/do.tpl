// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"{{if .table.Imports}}{{range $k,$v := .table.Imports}}
	"{{$k}}"{{end}}{{end}}
)

// {{.table.NameCaseCamel}} is the golang structure of table {{.table.Name}} for DAO operations like Where/Data.
type {{.table.NameCaseCamel}} struct {
	g.Meta `orm:"table:{{.table.Name}}, do:true"`{{range $i,$v := .table.Fields}}
	{{$v.NameCaseCamel}} interface{} // {{$v.Comment}}{{end}}
}
