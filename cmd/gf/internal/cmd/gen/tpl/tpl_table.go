package tpl

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/gogf/gf/v2/database/gdb"
)

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

// Table description
type Table struct {
	Name         string
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

// toTableFields description
//
// createTime: 2023-10-23 17:22:40
//
// author: hailaz
func (t *Table) toTableFields() {
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
func (t *Table) FieldsJsonStr() string {
	mapStr := make(map[string]interface{}, len(t.Fields))
	for _, v := range t.Fields {
		mapStr[v.Name] = v.Default
	}
	b, err := json.Marshal(mapStr)
	if err != nil {
		return ""
	}
	return string(b)
}
