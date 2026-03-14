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

// MSSQLParser implements SQLParser for SQL Server (T-SQL) DDL.
type MSSQLParser struct{}

// ParseCreateTable parses a single MSSQL CREATE TABLE statement.
func (p *MSSQLParser) ParseCreateTable(stmt string) (string, map[string]*gdb.TableField, error) {
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

// ParseAlterTable parses MSSQL ALTER TABLE statements.
func (p *MSSQLParser) ParseAlterTable(stmt string, tables map[string]map[string]*gdb.TableField) error {
	return parseAlterTableCommon(stmt, tables, p.parseColumnDef)
}

// ParseComment parses EXEC sp_addextendedproperty to extract column comments.
func (p *MSSQLParser) ParseComment(stmt string, tables map[string]map[string]*gdb.TableField) {
	upper := strings.ToUpper(strings.TrimSpace(stmt))
	if !strings.Contains(upper, "SP_ADDEXTENDEDPROPERTY") ||
		!strings.Contains(upper, "MS_DESCRIPTION") {
		return
	}

	// Extract quoted string values
	var values []string
	inQuote := false
	var current strings.Builder
	for i := 0; i < len(stmt); i++ {
		ch := stmt[i]
		if ch == '\'' {
			if inQuote {
				if i+1 < len(stmt) && stmt[i+1] == '\'' {
					current.WriteByte('\'')
					i++
					continue
				}
				values = append(values, current.String())
				current.Reset()
				inQuote = false
			} else {
				inQuote = true
			}
		} else if inQuote {
			current.WriteByte(ch)
		}
	}

	if len(values) < 8 {
		return
	}

	var (
		comment    string
		tableName  string
		columnName string
	)

	for i := 0; i < len(values)-1; i++ {
		switch strings.ToUpper(values[i]) {
		case "MS_DESCRIPTION":
			comment = values[i+1]
		case "TABLE":
			tableName = values[i+1]
		case "COLUMN":
			columnName = values[i+1]
		}
	}

	if tableName != "" && columnName != "" && comment != "" {
		if fields, ok := tables[tableName]; ok {
			if field, ok := fields[columnName]; ok {
				field.Comment = comment
			}
		}
	}
}

// parseColumnDef parses a single MSSQL column definition string into a TableField.
// It handles MSSQL-specific syntax including bracket-quoted identifiers and
// type parameters like varchar(max).
func (p *MSSQLParser) parseColumnDef(def string, index int) (*gdb.TableField, error) {
	tokens := mysqlTokenize(def)
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid column definition: %s", def)
	}

	field := &gdb.TableField{
		Index: index,
		Name:  unquoteIdentifier(tokens[0]),
		Null:  true,
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

// parseColumnAttributes parses MSSQL column constraint keywords including
// NOT NULL, NULL, PRIMARY KEY, UNIQUE, IDENTITY (auto-increment), and DEFAULT.
func (p *MSSQLParser) parseColumnAttributes(field *gdb.TableField, attrs string) {
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
				i++
			}
		case "UNIQUE":
			if field.Key == "" {
				field.Key = "UNI"
			}
		case "IDENTITY":
			field.Extra = "auto_increment"
			if i+1 < len(words) && strings.HasPrefix(words[i+1], "(") {
				i++
			}
		default:
			if strings.HasPrefix(upperWords[i], "IDENTITY(") || strings.HasPrefix(upperWords[i], "IDENTITY (") {
				field.Extra = "auto_increment"
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
