// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

// TableFields retrieves and returns the fields' information of specified table of current schema.
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result     gdb.Result
		link       gdb.Link
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, fmt.Sprintf(`PRAGMA TABLE_INFO(%s)`, d.QuoteWord(table)))
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		mKey := ""
		if m["pk"].Bool() {
			mKey = "pri"
		}
		fields[m["name"].String()] = &gdb.TableField{
			Index:   i,
			Name:    m["name"].String(),
			Type:    m["type"].String(),
			Key:     mKey,
			Default: m["dflt_value"].Val(),
			Null:    !m["notnull"].Bool(),
		}
	}
	return fields, nil
}
