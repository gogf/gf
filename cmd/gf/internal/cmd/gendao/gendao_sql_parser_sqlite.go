// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// SQLiteParser implements SQLParser for SQLite DDL.
type SQLiteParser struct{}

// ParseCreateTable parses a single SQLite CREATE TABLE statement.
func (p *SQLiteParser) ParseCreateTable(stmt string) (string, map[string]*gdb.TableField, error) {
	body, _, ok := extractBodyAndTrailing(stmt)
	if !ok {
		return "", nil, nil
	}

	parenIdx := strings.Index(stmt, "(")
	header := stmt[:parenIdx]
	tableName := extractTableName(header)
	if tableName == "" {
		return "", nil, fmt.Errorf("cannot extract table name from: %s", header)
	}

	columnDefs := splitColumns(body)
	fields := make(map[string]*gdb.TableField)
	pkColumns := findPrimaryKeysFromConstraints(columnDefs)

	fieldIndex := 0
	for _, def := range columnDefs {
		def = strings.TrimSpace(def)
		if def == "" {
			continue
		}
		firstWord := strings.ToUpper(strings.Fields(def)[0])
		if isConstraintKeyword(firstWord) {
			continue
		}

		field, err := p.parseColumnDef(def, fieldIndex)
		if err != nil {
			continue
		}
		if field != nil {
			fields[field.Name] = field
			fieldIndex++
		}
	}

	for _, pkCol := range pkColumns {
		if f, ok := fields[pkCol]; ok {
			f.Key = "PRI"
		}
	}

	return tableName, fields, nil
}

// ParseAlterTable parses SQLite ALTER TABLE statements.
// Note: SQLite only supports ADD COLUMN and RENAME COLUMN in ALTER TABLE.
func (p *SQLiteParser) ParseAlterTable(stmt string, tables map[string]map[string]*gdb.TableField) error {
	return parseAlterTableCommon(stmt, tables, p.parseColumnDef)
}

// ParseComment is a no-op for SQLite as it doesn't support COMMENT ON statements.
func (p *SQLiteParser) ParseComment(stmt string, tables map[string]map[string]*gdb.TableField) {
	// SQLite does not support comments on columns.
}

// parseColumnDef parses a single SQLite column definition string into a TableField.
// SQLite has flexible typing (type affinity), so columns may have no explicit type,
// in which case "text" is used as the default type.
func (p *SQLiteParser) parseColumnDef(def string, index int) (*gdb.TableField, error) {
	tokens := mysqlTokenize(def)
	if len(tokens) < 1 {
		return nil, fmt.Errorf("invalid column definition: %s", def)
	}

	field := &gdb.TableField{
		Index: index,
		Name:  unquoteIdentifier(tokens[0]),
		Null:  true,
	}

	if len(tokens) < 2 {
		field.Type = "text"
		return field, nil
	}

	field.Type = tokens[1]

	rest := ""
	if len(tokens) > 2 {
		rest = strings.Join(tokens[2:], " ")
	}

	if !strings.Contains(field.Type, "(") && strings.HasPrefix(strings.TrimSpace(rest), "(") {
		end := strings.Index(rest, ")")
		if end >= 0 {
			field.Type += rest[:end+1]
			rest = strings.TrimSpace(rest[end+1:])
		}
	}

	p.parseColumnAttributes(field, rest)

	return field, nil
}

// parseColumnAttributes parses SQLite column constraint keywords including
// NOT NULL, NULL, PRIMARY KEY (with optional AUTOINCREMENT), UNIQUE, and DEFAULT.
func (p *SQLiteParser) parseColumnAttributes(field *gdb.TableField, attrs string) {
	words := strings.Fields(attrs)
	upperWords := strings.Fields(strings.ToUpper(attrs))

	for i := 0; i < len(upperWords); i++ {
		switch upperWords[i] {
		case "NOT":
			if i+1 < len(upperWords) && upperWords[i+1] == "NULL" {
				field.Null = false
				i++
			}
		case "NULL":
			field.Null = true
		case "PRIMARY":
			if i+1 < len(upperWords) && upperWords[i+1] == "KEY" {
				field.Key = "PRI"
				field.Null = false
				i++
				if i+1 < len(upperWords) && upperWords[i+1] == "AUTOINCREMENT" {
					field.Extra = "auto_increment"
					i++
				}
			}
		case "AUTOINCREMENT":
			field.Extra = "auto_increment"
		case "UNIQUE":
			if field.Key == "" {
				field.Key = "UNI"
			}
		case "DEFAULT":
			if i+1 < len(words) {
				defaultVal, _ := extractDefaultValue("DEFAULT " + strings.Join(words[i+1:], " "))
				field.Default = defaultVal
				if defaultVal != nil {
					i++
				}
			}
		}
	}
}
