// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

// TableFields retrieves and returns the fields' information of specified table of current
// schema.
//
// The parameter `link` is optional, if given nil it automatically retrieves a raw sql connection
// as its link to proceed necessary sql query.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the TableField struct has a "Index" field marks its sequence in
// the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the
// process restarts.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result             gdb.Result
		link               gdb.Link
		usedSchema         = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		mariaJsonFiledName = gset.New(false)
	)
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	result, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(`SHOW FULL COLUMNS FROM %s`, d.QuoteWord(table)),
	)
	if err != nil {
		return nil, err
	}

	// Compatible with mariaDB json type
	dbType := d.GetConfig().Type
	if dbType == "mariadb" {
		var checkConstraintResult gdb.Result
		checkConstraintResult, err = d.DoSelect(
			ctx, link,
			fmt.Sprintf(`SELECT CONSTRAINT_NAME as filedName, CHECK_CLAUSE as filedCheck FROM information_schema.CHECK_CONSTRAINTS WHERE TABLE_NAME = '%s'`, table),
		)
		if err != nil {
			return nil, err
		}
		for _, m := range checkConstraintResult {
			if gstr.HasPrefix(m["filedCheck"].String(), "json_valid") {
				mariaJsonFiledName.Add(m["filedName"].String())
			}
		}
	}

	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		// if filed exists in mariaJsonFiledName, replace its Type with "json"
		if mariaJsonFiledName.Size() != 0 && mariaJsonFiledName.Contains(m["Field"].String()) {
			m["Type"] = gvar.New("json")
		}
		fields[m["Field"].String()] = &gdb.TableField{
			Index:   i,
			Name:    m["Field"].String(),
			Type:    m["Type"].String(),
			Null:    m["Null"].Bool(),
			Key:     m["Key"].String(),
			Default: m["Default"].Val(),
			Extra:   m["Extra"].String(),
			Comment: m["Comment"].String(),
		}
	}
	return fields, nil
}
