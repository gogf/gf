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

// OracleParser implements SQLParser for Oracle/DM DDL.
type OracleParser struct{}

// ParseCreateTable parses a single Oracle CREATE TABLE statement.
func (p *OracleParser) ParseCreateTable(stmt string) (string, map[string]*gdb.TableField, error) {
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
		upperPk := strings.ToUpper(pkCol)
		if f, ok := fields[upperPk]; ok {
			f.Key = "PRI"
		}
	}

	return tableName, fields, nil
}

// ParseAlterTable parses Oracle ALTER TABLE statements.
func (p *OracleParser) ParseAlterTable(stmt string, tables map[string]map[string]*gdb.TableField) error {
	return parseAlterTableCommon(stmt, tables, p.parseColumnDef)
}

// ParseComment parses COMMENT ON COLUMN table.column IS 'comment'.
func (p *OracleParser) ParseComment(stmt string, tables map[string]map[string]*gdb.TableField) {
	upper := strings.ToUpper(strings.TrimSpace(stmt))
	if !strings.HasPrefix(upper, "COMMENT ON COLUMN") {
		return
	}

	rest := strings.TrimSpace(stmt[len("COMMENT ON COLUMN"):])
	isIdx := strings.Index(strings.ToUpper(rest), " IS ")
	if isIdx < 0 {
		return
	}
	ref := strings.TrimSpace(rest[:isIdx])
	comment := strings.TrimSpace(rest[isIdx+4:])

	if len(comment) >= 2 && comment[0] == '\'' && comment[len(comment)-1] == '\'' {
		comment = comment[1 : len(comment)-1]
		comment = strings.ReplaceAll(comment, "''", "'")
	}

	parts := strings.Split(ref, ".")
	var tableName, columnName string
	switch len(parts) {
	case 2:
		tableName = unquoteIdentifier(parts[0])
		columnName = unquoteIdentifier(parts[1])
	case 3:
		tableName = unquoteIdentifier(parts[1])
		columnName = unquoteIdentifier(parts[2])
	default:
		return
	}

	if fields, ok := tables[tableName]; ok {
		if field, ok := fields[columnName]; ok {
			field.Comment = comment
		}
	}
}

// parseColumnDef parses a single Oracle column definition string into a TableField.
// It handles Oracle-specific types including TIMESTAMP WITH TIME ZONE and
// TIMESTAMP WITH LOCAL TIME ZONE.
func (p *OracleParser) parseColumnDef(def string, index int) (*gdb.TableField, error) {
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

	// Handle TIMESTAMP WITH TIME ZONE / WITH LOCAL TIME ZONE
	upperType := strings.ToUpper(field.Type)
	upperRest := strings.ToUpper(rest)
	if upperType == "TIMESTAMP" {
		if strings.HasPrefix(upperRest, "WITH LOCAL TIME ZONE") {
			field.Type = "timestamp with local time zone"
			rest = strings.TrimSpace(rest[len("WITH LOCAL TIME ZONE"):])
		} else if strings.HasPrefix(upperRest, "WITH TIME ZONE") {
			field.Type = "timestamp with time zone"
			rest = strings.TrimSpace(rest[len("WITH TIME ZONE"):])
		}
	}

	p.parseColumnAttributes(field, rest)

	return field, nil
}

// parseColumnAttributes parses Oracle column constraint keywords including
// NOT NULL, NULL, PRIMARY KEY, UNIQUE, DEFAULT, and GENERATED ... AS IDENTITY.
func (p *OracleParser) parseColumnAttributes(field *gdb.TableField, attrs string) {
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
		case "DEFAULT":
			if i+1 < len(words) {
				defaultVal, _ := extractDefaultValue("DEFAULT " + strings.Join(words[i+1:], " "))
				field.Default = defaultVal
				if defaultVal != nil {
					i++
				}
			}
		case "GENERATED":
			rest := strings.Join(upperWords[i:], " ")
			if strings.Contains(rest, "AS IDENTITY") {
				field.Extra = "auto_increment"
				for j := i + 1; j < len(upperWords); j++ {
					if upperWords[j] == "IDENTITY" {
						i = j
						break
					}
				}
			}
		}
	}
}
