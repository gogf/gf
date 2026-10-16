// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

// GetDB returns the underlying DB.
func (c *Core) GetDB() DB {
	return c.db
}

// GetLink creates and returns the underlying database link object with transaction checks.
// The parameter `master` specifies whether using the master node if master-slave configured.
func (c *Core) GetLink(ctx context.Context, master bool, schema string) (Link, error) {
	tx := TXFromCtx(ctx, c.db.GetGroup())
	if tx != nil {
		return &txLink{tx.GetSqlTX()}, nil
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
	if !gregex.IsMatchString(regularFieldNameWithCommaRegPattern, s) {
		return s
	}
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
	return doQuoteTableName(table, c.db.GetPrefix(), charLeft, charRight)
}

// GetChars returns the security char for current database.
// It does nothing in default.
func (c *Core) GetChars() (charLeft string, charRight string) {
	return "", ""
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (c *Core) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	return
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
func (c *Core) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	return
}

// ClearTableFields clears the cached field data for the specified table so that the
// next call to TableFields re-queries the database. The table's existence marker is
// preserved so that HasTable continues to return true without a DB round-trip.
func (c *Core) ClearTableFields(ctx context.Context, table string, schema ...string) (err error) {
	c.db.GetCore().registry.Set(
		c.db.GetGroup(),
		gutil.GetOrDefaultStr(c.db.GetSchema(), schema...),
		table,
		nil,
	)
	return
}

// ClearTableFieldsAll removes all cached table fields of current configuration group.
func (c *Core) ClearTableFieldsAll(ctx context.Context) (err error) {
	c.db.GetCore().registry.ClearAll()
	return
}

// ClearCache removes cached sql result of certain table.
func (c *Core) ClearCache(ctx context.Context, table string) (err error) {
	var (
		keys, _     = c.db.GetCache().KeyStrings(ctx)
		cachePrefix = fmt.Sprintf(`%s%s@`, cachePrefixSelectCache, table)
		removedKeys = make([]any, 0)
	)
	for _, key := range keys {
		if gstr.HasPrefix(key, cachePrefix) {
			removedKeys = append(removedKeys, key)
		}
	}
	if len(removedKeys) > 0 {
		err = c.db.GetCache().Removes(ctx, removedKeys)
	}
	return
}

// ClearCacheAll removes all cached sql result from cache
func (c *Core) ClearCacheAll(ctx context.Context) (err error) {
	if err = c.db.GetCache().Clear(ctx); err != nil {
		return err
	}
	c.db.GetCore().registry.ClearAll()
	return
}

// HasField determine whether the field exists in the table.
func (c *Core) HasField(ctx context.Context, table, field string, schema ...string) (bool, error) {
	table = c.db.GetTableNameForFields(table)
	tableFields, err := c.db.TableFields(ctx, table, schema...)
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
	foldedField := c.db.FoldIdentifier(field)
	for _, f := range fieldsArray {
		if f == field || f == foldedField {
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
		guessedTableName string
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

// tableReferenceModifiers are the keywords that may precede the table of a table reference.
var tableReferenceModifiers = map[string]bool{"ONLY": true, "LATERAL": true}

// tableNameForFields returns the table name that the fields of the table expression `table` of a
// statement are looked up by. It is the table of a single table with an alias, like "user u",
// "db.user AS u" or "`my user` `u`", as it is written without the alias, so that it is looked up
// as the same table without an alias. Any other table expression is returned as it is.
func (c *Core) tableNameForFields(table string) string {
	words, ok := splitSqlWords(table)
	hasAlias := len(words) == 2 && !strings.EqualFold(words[1], "AS") ||
		len(words) == 3 && strings.EqualFold(words[1], "AS")
	if !ok || !hasAlias {
		return table
	}
	alias := words[len(words)-1]
	if !isSqlIdentifier(words[0], true) || tableReferenceModifiers[strings.ToUpper(words[0])] ||
		!isSqlIdentifier(alias, false) {
		return table
	}
	return words[0]
}

// splitSqlWords splits `s` into its words separated by whitespace, taking an identifier quoted
// with '`' or '"' as part of a word whatever it contains. It returns false if a quote is not
// closed or if `s` has a comma outside quotes.
func splitSqlWords(s string) (words []string, ok bool) {
	var begin = -1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r':
			if begin >= 0 {
				words = append(words, s[begin:i])
				begin = -1
			}
			continue
		case ',':
			return nil, false
		case '`', '"':
			end := closingSqlQuote(s, i)
			if end < 0 {
				return nil, false
			}
			if begin < 0 {
				begin = i
			}
			i = end
			continue
		}
		if begin < 0 {
			begin = i
		}
	}
	if begin >= 0 {
		words = append(words, s[begin:])
	}
	return words, true
}

// closingSqlQuote returns the index of the quote closing the quoted identifier that begins at
// index `open` of `s`, taking a doubled quote as an escaped one, or -1 if it is not closed.
func closingSqlQuote(s string, open int) int {
	for i := open + 1; i < len(s); i++ {
		if s[i] != s[open] {
			continue
		}
		if i+1 < len(s) && s[i+1] == s[open] {
			i++
			continue
		}
		return i
	}
	return -1
}

// isSqlIdentifier reports whether `s` is an identifier, unquoted or quoted with '`' or '"', or,
// if `qualified` is true, identifiers joined by '.' like "db.user" or "`db`.`user`".
func isSqlIdentifier(s string, qualified bool) bool {
	for i := 0; ; i++ {
		if i >= len(s) {
			return false
		}
		if s[i] == '`' || s[i] == '"' {
			end := closingSqlQuote(s, i)
			if end <= i+1 {
				return false
			}
			i = end + 1
		} else {
			begin := i
			for i < len(s) && isSqlIdentifierChar(s[i]) {
				i++
			}
			if i == begin {
				return false
			}
		}
		if i == len(s) {
			return true
		}
		if s[i] != '.' || !qualified {
			return false
		}
	}
}

// isSqlIdentifierChar reports whether `c` can be part of an unquoted identifier: a letter, a digit,
// '_', '$', '#', or a byte of a multibyte character.
func isSqlIdentifierChar(c byte) bool {
	return c == '_' || c == '$' || c == '#' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

// GetPrimaryKeys retrieves and returns the primary key field names of the specified table.
// This method extracts primary key information from TableFields.
// The parameter `schema` is optional, if not specified it uses the default schema.
func (c *Core) GetPrimaryKeys(ctx context.Context, table string, schema ...string) ([]string, error) {
	tableFields, err := c.db.TableFields(ctx, c.tableNameForFields(table), schema...)
	if err != nil {
		return nil, err
	}

	var primaryKeys []string
	for _, field := range tableFields {
		if strings.EqualFold(field.Key, "pri") {
			primaryKeys = append(primaryKeys, field.Name)
		}
	}

	return primaryKeys, nil
}
