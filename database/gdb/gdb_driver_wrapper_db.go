// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// DriverWrapperDB is a DB wrapper for extending features with embedded DB.
type DriverWrapperDB struct {
	DB
}

// Open creates and returns an underlying sql.DB object for pgsql.
// https://pkg.go.dev/github.com/lib/pq
func (d *DriverWrapperDB) Open(config *ConfigNode) (db *sql.DB, err error) {
	if config.Link != "" {
		config = parseConfigNodeLink(config)
	}
	return d.DB.Open(config)
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverWrapperDB) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	ctx = context.WithValue(ctx, ctxKeyInternalProducedSQL, struct{}{})
	return d.DB.Tables(ctx, schema...)
}

// TableFields retrieves and returns the fields' information of specified table of current
// schema.
//
// The parameter `link` is optional, if given nil it automatically retrieves a raw sql connection
// as its link to proceed necessary sql query.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the TableField struct has an "Index" field marks its sequence in
// the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the
// process restarts.
func (d *DriverWrapperDB) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(
			gcode.CodeInvalidParameter,
			"function TableFields supports only single table operations",
		)
	}
	var (
		cacheKey = fmt.Sprintf(
			`%s%s@%s#%s`,
			cachePrefixTableFields,
			d.GetGroup(),
			gutil.GetOrDefaultStr(d.GetSchema(), schema...),
			table,
		)
		value = tableFieldsMap.GetOrSetFuncLock(cacheKey, func() interface{} {
			ctx = context.WithValue(ctx, ctxKeyInternalProducedSQL, struct{}{})
			fields, err = d.DB.TableFields(ctx, table, schema...)
			if err != nil {
				return nil
			}
			return fields
		})
	)
	if value != nil {
		fields = value.(map[string]*TableField)
	}
	return
}

func parseConfigNodeLink(node *ConfigNode) *ConfigNode {
	var match []string
	// It firstly parses `link` using with type pattern.
	match, _ = gregex.MatchString(linkPatternWithType, node.Link)
	if len(match) > 6 {
		node.Type = match[1]
		node.User = match[2]
		node.Pass = match[3]
		node.Protocol = match[4]
		node.Host = match[5]
		node.Port = match[6]
		node.Name = match[7]
		if len(match) > 7 {
			node.Extra = match[8]
		}
		node.Link = ""
	} else {
		// Else it parses `link` using without type pattern.
		match, _ = gregex.MatchString(linkPatternWithoutType, node.Link)
		if len(match) > 6 {
			node.User = match[1]
			node.Pass = match[2]
			node.Protocol = match[3]
			node.Host = match[4]
			node.Port = match[5]
			node.Name = match[6]
			if len(match) > 7 {
				node.Extra = match[7]
			}
			node.Link = ""
		}
	}
	if node.Extra != "" {
		if m, _ := gstr.Parse(node.Extra); len(m) > 0 {
			_ = gconv.Struct(m, &node)
		}
	}
	// Default value checks.
	if node.Charset == "" {
		node.Charset = defaultCharset
	}
	if node.Protocol == "" {
		node.Protocol = defaultProtocol
	}
	return node
}
