// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// WithDB injects given db object into context and returns a new context.
func WithDB(ctx context.Context, db DB) context.Context {
	if db == nil {
		return ctx
	}
	dbCtx := db.GetCtx()
	if ctxDb := DBFromCtx(dbCtx); ctxDb != nil {
		return dbCtx
	}
	ctx = context.WithValue(ctx, contextKeyForDB, db)
	return ctx
}

// DBFromCtx retrieves and returns DB object from context.
func DBFromCtx(ctx context.Context) DB {
	if ctx == nil {
		return nil
	}
	v := ctx.Value(contextKeyForDB)
	if v != nil {
		return v.(DB)
	}
	return nil
}

// GetLink creates and returns the underlying database link object with transaction checks.
// The parameter `master` specifies whether using the master node if master-slave configured.
func (c *Core) GetLink(ctx context.Context, master bool, schema string) (Link, error) {
	tx := TXFromCtx(ctx, c.db.GetGroup())
	if tx != nil {
		return &txLink{tx.tx}, nil
	}
	if master {
		link, err := c.db.GetCore().MasterLink(schema)
		if err != nil {
			return nil, err
		}
		return link, nil
	}
	link, err := c.db.GetCore().SlaveLink(schema)
	if err != nil {
		return nil, err
	}
	return link, nil
}

// MasterLink acts like function Master but with additional `schema` parameter specifying
// the schema for the connection. It is defined for internal usage.
// Also see Master.
func (c *Core) MasterLink(schema ...string) (Link, error) {
	db, err := c.db.Master(schema...)
	if err != nil {
		return nil, err
	}
	return &dbLink{
		DB:         db,
		isOnMaster: true,
	}, nil
}

// SlaveLink acts like function Slave but with additional `schema` parameter specifying
// the schema for the connection. It is defined for internal usage.
// Also see Slave.
func (c *Core) SlaveLink(schema ...string) (Link, error) {
	db, err := c.db.Slave(schema...)
	if err != nil {
		return nil, err
	}
	return &dbLink{
		DB:         db,
		isOnMaster: false,
	}, nil
}

// QuoteWord checks given string `s` a word,
// if true it quotes `s` with security chars of the database
// and returns the quoted string; or else it returns `s` without any change.
//
// The meaning of a `word` can be considered as a column name.
func (c *Core) QuoteWord(s string) string {
	s = gstr.Trim(s)
	if s == "" {
		return s
	}
	charLeft, charRight := c.db.GetChars()
	return doQuoteWord(s, charLeft, charRight)
}

// QuoteString quotes string with quote chars. Strings like:
// "user", "user u", "user,user_detail", "user u, user_detail ut", "u.id asc".
//
// The meaning of a `string` can be considered as part of a statement string including columns.
func (c *Core) QuoteString(s string) string {
	charLeft, charRight := c.db.GetChars()
	return doQuoteString(s, charLeft, charRight)
}

// QuotePrefixTableName adds prefix string and quotes chars for the table.
// It handles table string like:
// "user", "user u",
// "user,user_detail",
// "user u, user_detail ut",
// "user as u, user_detail as ut".
//
// Note that, this will automatically checks the table prefix whether already added,
// if true it does nothing to the table name, or else adds the prefix to the table name.
func (c *Core) QuotePrefixTableName(table string) string {
	charLeft, charRight := c.db.GetChars()
	return doHandleTableName(table, c.db.GetPrefix(), charLeft, charRight)
}

// GetChars returns the security char for current database.
// It does nothing in default.
func (c *Core) GetChars() (charLeft string, charRight string) {
	return "", ""
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
//
// It does nothing in default.
func (c *Core) Tables(schema ...string) (tables []string, err error) {
	return
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
//
// Note that it returns a map containing the field name and its corresponding fields.
// As a map is unsorted, the TableField struct has an "Index" field marks its sequence in the fields.
//
// It's using cache feature to enhance the performance, which is never expired util the process restarts.
//
// It does nothing in default.
func (c *Core) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	// It does nothing if given table is empty, especially in sub-query.
	if table == "" {
		return map[string]*TableField{}, nil
	}
	return c.db.TableFields(c.GetCtx(), table, schema...)
}

// HasField determine whether the field exists in the table.
func (c *Core) HasField(table, field string, schema ...string) (bool, error) {
	table = c.guessPrimaryTableName(table)
	tableFields, err := c.TableFields(table, schema...)
	if err != nil {
		return false, err
	}
	if len(tableFields) == 0 {
		return false, gerror.NewCodef(
			gcode.CodeNotFound,
			`empty table fields for table "%s"`, table,
		)
	}
	fieldsArray := make([]string, len(tableFields))
	for k, v := range tableFields {
		fieldsArray[v.Index] = k
	}
	charLeft, charRight := c.db.GetChars()
	field = gstr.Trim(field, charLeft+charRight)
	for _, f := range fieldsArray {
		if f == field {
			return true, nil
		}
	}
	return false, nil
}

// guessPrimaryTableName parses and returns the primary table name.
func (c *Core) guessPrimaryTableName(tableStr string) string {
	if tableStr == "" {
		return ""
	}
	var (
		guessedTableName = ""
		array1           = gstr.SplitAndTrim(tableStr, ",")
		array2           = gstr.SplitAndTrim(array1[0], " ")
		array3           = gstr.SplitAndTrim(array2[0], ".")
	)
	if len(array3) >= 2 {
		guessedTableName = array3[1]
	} else {
		guessedTableName = array3[0]
	}
	charL, charR := c.db.GetChars()
	if charL != "" || charR != "" {
		guessedTableName = gstr.Trim(guessedTableName, charL+charR)
	}
	if !gregex.IsMatchString(regularFieldNameRegPattern, guessedTableName) {
		return ""
	}
	return guessedTableName
}
