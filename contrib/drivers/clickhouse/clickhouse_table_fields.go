// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	tableFieldsColumns = `name,position,default_expression,comment,type,is_in_partition_key,is_in_sorting_key,is_in_primary_key,is_in_sampling_key`
)

// TableFields retrieves and returns the fields' information of specified table of current schema.
// A table prefixed with its database, like "db.table", is read from that database.
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result     gdb.Result
		link       gdb.Link
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		database   string
	)
	if database, table = splitTableName(table); database != "" {
		usedSchema = database
	}
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	var (
		getColumnsSql = fmt.Sprintf(
			"select %s from `system`.columns c where `database` = currentDatabase() and `table` = ?",
			tableFieldsColumns,
		)
	)
	result, err = d.DoSelect(ctx, link, getColumnsSql, table)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for _, m := range result {
		var (
			isNull    = false
			fieldType = m["type"].String()
			key       string
		)
		if strings.HasPrefix(fieldType, "Nullable(") && strings.HasSuffix(fieldType, ")") {
			isNull = true
			fieldType = fieldType[len("Nullable(") : len(fieldType)-1]
		}
		if m["is_in_primary_key"].Bool() {
			key = "PRI"
		}
		position := m["position"].Int()
		if result[0]["position"].Int() != 0 {
			position -= 1
		}
		fields[m["name"].String()] = &gdb.TableField{
			Index:   position,
			Name:    m["name"].String(),
			Default: m["default_expression"].Val(),
			Comment: m["comment"].String(),
			Key:     key,
			Type:    fieldType,
			Null:    isNull,
		}
	}
	return fields, nil
}

// splitTableName returns the database and the name of the table `table`, unquoted, like "db" and
// "t" for "db.t" or "`db`.`t`", with an empty database if `table` is not prefixed with one. An
// alias following the name is ignored.
func splitTableName(table string) (database, name string) {
	table = strings.TrimSpace(table)
	end := identifierEnd(table, 0)
	if end < len(table) && table[end] == '.' {
		return unquoteIdentifier(table[:end]), unquoteIdentifier(table[end+1 : identifierEnd(table, end+1)])
	}
	return "", unquoteIdentifier(table[:end])
}

// identifierEnd returns the index following the identifier that begins at index `i` of `s`,
// which is either quoted with backticks or double quotes or ends before a '.' or a white space.
func identifierEnd(s string, i int) int {
	if i < len(s) && (s[i] == '`' || s[i] == '"') {
		for j := i + 1; j < len(s); j++ {
			switch s[j] {
			case '\\':
				j++
			case s[i]:
				if j+1 < len(s) && s[j+1] == s[i] {
					j++
					continue
				}
				return j + 1
			}
		}
		return len(s)
	}
	for i < len(s) && s[i] != '.' && s[i] != ' ' && s[i] != '\t' && s[i] != '\r' && s[i] != '\n' {
		i++
	}
	return i
}

// unquoteIdentifier returns the identifier `s` without its quotes and with its escaped
// characters unescaped, or `s` itself if it is not quoted.
func unquoteIdentifier(s string) string {
	if len(s) < 2 || (s[0] != '`' && s[0] != '"') || s[len(s)-1] != s[0] {
		return s
	}
	var (
		quote = s[0]
		b     strings.Builder
	)
	s = s[1 : len(s)-1]
	for i := 0; i < len(s); i++ {
		if i+1 < len(s) && (s[i] == '\\' || (s[i] == quote && s[i+1] == quote)) {
			i++
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
