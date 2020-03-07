// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import "database/sql"

// GetMaster acts like function Master but with additional <schema> parameter specifying
// the schema for the connection. It is defined for internal usage.
// Also see Master.
func (c *Core) GetMaster(schema ...string) (*sql.DB, error) {
	return c.getSqlDb(true, schema...)
}

// GetSlave acts like function Slave but with additional <schema> parameter specifying
// the schema for the connection. It is defined for internal usage.
// Also see Slave.
func (c *Core) GetSlave(schema ...string) (*sql.DB, error) {
	return c.getSqlDb(false, schema...)
}

// QuoteWord checks given string <s> a word, if true quotes it with security chars of the database
// and returns the quoted string; or else return <s> without any change.
func (c *Core) QuoteWord(s string) string {
	charLeft, charRight := c.DB.GetChars()
	return doQuoteWord(s, charLeft, charRight)
}

// QuoteString quotes string with quote chars. Strings like:
// "user", "user u", "user,user_detail", "user u, user_detail ut", "u.id asc".
func (c *Core) QuoteString(s string) string {
	charLeft, charRight := c.DB.GetChars()
	return doQuoteString(s, charLeft, charRight)
}

// GetChars returns the security char for current database.
// It does nothing in default.
func (c *Core) GetChars() (charLeft string, charRight string) {
	return "", ""
}

// HandleSqlBeforeExec handles the sql before posts it to database.
// It does nothing in default.
func (c *Core) HandleSqlBeforeExec(sql string) string {
	return sql
}
