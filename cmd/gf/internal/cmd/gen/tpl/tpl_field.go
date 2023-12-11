package tpl

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
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
