// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity
{{if .table.Imports}}
import ({{range $k,$v := .table.Imports}}
	"{{$k}}"{{end}}
)
{{end}}
// {{.table.NameCaseCamel}} is the golang structure for table {{.table.Name}}.
type {{.table.NameCaseCamel}} struct { {{range $i,$v := .table.Fields}}
	{{$v.NameCaseCamel}} {{$v.LocalType}} `json:"{{$v.NameJsonCase}}" orm:"{{$v.Name}}" description:"{{$v.Comment}}"` // {{$v.Comment}}{{end}}
}
