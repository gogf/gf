package tpl

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
)

// Table description
type Table struct {
	Name         string
	PackageName  string
	db           gdb.DB
	Fields       TableFields
	FieldsSource map[string]*gdb.TableField
	Imports      map[string]struct{}
}

type Tables []*Table

// NewTable description
//
// createTime: 2023-12-11 16:17:33
//
// author: hailaz
func NewTable(ctx context.Context, db gdb.DB, tableName string) (*Table, error) {
	fields, err := db.TableFields(ctx, tableName)
	if err != nil {
		return nil, err
	}
	table := Table{
		Name:         tableName,
		FieldsSource: fields,
		db:           db,
		Imports:      make(map[string]struct{}),
	}
	table.toTableFields()
	return &table, nil
}

// Name description
//
// createTime: 2023-10-23 16:17:30
//
// author: hailaz
func (t *Table) Show() string {
	return fmt.Sprintf("table name is %s", t.Name)
}

// CaseCamel description
func (t *Table) CaseCamel() string {
	return gstr.CaseCamel(t.Name)
}

// CaseCamelLower description
func (t *Table) CaseCamelLower() string {
	return gstr.CaseCamelLower(t.Name)
}

// CaseSnake description
func (t *Table) CaseSnake() string {
	return gstr.CaseSnake(t.Name)
}

// CaseKebabScreaming description
func (t *Table) CaseKebabScreaming() string {
	return gstr.CaseKebabScreaming(t.Name)
}

// toTableFields description
//
// createTime: 2023-10-23 17:22:40
//
// author: hailaz
func (t *Table) toTableFields() {
	if len(t.Fields) > 0 {
		return
	}
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
func (t *Table) SortFields(isReverse bool) {
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
func (t *Table) FieldsJsonStr(caseName string) string {
	mapStr := make(map[string]interface{}, len(t.Fields))
	for _, v := range t.Fields {
		mapStr[v.NameCase(caseName)] = v.Default
	}
	b, err := json.MarshalIndent(mapStr, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}
