// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// DoFilter handles the sql before posts it to database.
// This method helps handle MySQL-specific issues including key length limitations.
// 
// For MySQL tables using utf8mb4 charset, this method automatically adds ROW_FORMAT=DYNAMIC
// to CREATE TABLE statements to prevent "Specified key was too long; max key length is 1000 bytes" errors.
// This is particularly important for compatibility when upgrading from older GoFrame versions.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []any,
) (newSql string, newArgs []any, err error) {
	newSql, newArgs, err = d.Core.DoFilter(ctx, link, sql, args)
	if err != nil {
		return newSql, newArgs, err
	}
	
	// Handle MySQL-specific SQL filtering to prevent key length issues
	// This is particularly important for compatibility between GoFrame versions
	newSql = d.handleMySQLKeyLengthCompatibility(newSql)
	
	return newSql, newArgs, err
}

// handleMySQLKeyLengthCompatibility modifies SQL statements to be more compatible
// with MySQL key length limitations, especially for upgrade scenarios from older GoFrame versions
func (d *Driver) handleMySQLKeyLengthCompatibility(sql string) string {
	// For CREATE TABLE statements with utf8mb4 charset, ensure key length compatibility
	// This helps prevent "Specified key was too long; max key length is 1000 bytes" errors
	sqlUpper := strings.ToUpper(sql)
	sqlLower := strings.ToLower(sql)
	
	if strings.Contains(sqlUpper, "CREATE TABLE") && 
	   (strings.Contains(sqlLower, "utf8mb4") || strings.Contains(sqlLower, "charset=utf8mb4")) {
		// Add ROW_FORMAT=DYNAMIC to enable larger key prefixes when using utf8mb4
		if !strings.Contains(sqlUpper, "ROW_FORMAT") {
			// Insert ROW_FORMAT=DYNAMIC before ENGINE clause if it exists
			if strings.Contains(sqlUpper, "ENGINE=") {
				sql = strings.Replace(sql, "ENGINE=", "ROW_FORMAT=DYNAMIC ENGINE=", 1)
			} else if strings.Contains(sqlUpper, "ENGINE ") {
				// Handle case where there's a space after ENGINE
				sql = strings.Replace(sql, "ENGINE ", "ROW_FORMAT=DYNAMIC ENGINE ", 1)
			} else {
				// Append ROW_FORMAT=DYNAMIC at the end of CREATE TABLE statement
				sql = strings.TrimSuffix(strings.TrimSpace(sql), ";")
				sql += " ROW_FORMAT=DYNAMIC"
				if !strings.HasSuffix(sql, ";") {
					sql += ";"
				}
			}
		}
	}
	
	return sql
}
