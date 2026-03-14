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

// MySQLParser implements SQLParser for MySQL/MariaDB/TiDB DDL.
type MySQLParser struct{}

// ParseCreateTable parses a single MySQL CREATE TABLE statement.
func (p *MySQLParser) ParseCreateTable(stmt string) (string, map[string]*gdb.TableField, error) {
	body, trailing, ok := extractBodyAndTrailing(stmt)
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

	// Extract inline comments from trailing table options (not used for field generation)
	_ = trailing

	return tableName, fields, nil
}

// ParseAlterTable parses MySQL ALTER TABLE statements.
func (p *MySQLParser) ParseAlterTable(stmt string, tables map[string]map[string]*gdb.TableField) error {
	return parseAlterTableCommon(stmt, tables, p.parseColumnDef)
}

// ParseComment handles MySQL-style comments (inline COMMENT keyword is handled in parseColumnDef).
func (p *MySQLParser) ParseComment(stmt string, tables map[string]map[string]*gdb.TableField) {
	// MySQL uses inline COMMENT 'xxx' in column definitions,
	// which is already handled by parseColumnDef. No separate COMMENT ON statement.
}

// parseColumnDef parses a single MySQL column definition string into a TableField.
// It extracts the column name, data type (including UNSIGNED modifier), and delegates
// attribute parsing (NULL, DEFAULT, PRIMARY KEY, COMMENT, etc.) to parseColumnAttributes.
func (p *MySQLParser) parseColumnDef(def string, index int) (*gdb.TableField, error) {
	tokens := mysqlTokenize(def)
	if len(tokens) < 2 {
		return nil, fmt.Errorf("invalid column definition: %s", def)
	}

	field := &gdb.TableField{
		Index: index,
		Name:  unquoteIdentifier(tokens[0]),
		Null:  true,
	}

	typeStr := tokens[1]
	rest := ""
	if len(tokens) > 2 {
		rest = strings.Join(tokens[2:], " ")
	}

	// Check if rest starts with '(' meaning the type params are in rest
	if !strings.Contains(typeStr, "(") && strings.HasPrefix(strings.TrimSpace(rest), "(") {
		endParen := strings.Index(rest, ")")
		if endParen >= 0 {
			typeStr += rest[:endParen+1]
			rest = strings.TrimSpace(rest[endParen+1:])
		}
	}

	field.Type = typeStr

	// Handle UNSIGNED
	upperRest := strings.ToUpper(rest)
	if strings.HasPrefix(upperRest, "UNSIGNED") {
		field.Type += " unsigned"
		rest = strings.TrimSpace(rest[8:])
	}

	p.parseColumnAttributes(field, rest)

	return field, nil
}

// parseColumnAttributes parses MySQL column constraint keywords from the attribute string
// following the column type. It handles NOT NULL, NULL, PRIMARY KEY, UNIQUE, AUTO_INCREMENT,
// DEFAULT, COMMENT, and ON UPDATE clauses.
func (p *MySQLParser) parseColumnAttributes(field *gdb.TableField, attrs string) {
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
			if i+1 < len(upperWords) && upperWords[i+1] == "KEY" {
				i++
			}
		case "KEY":
			if field.Key == "" {
				field.Key = "MUL"
			}
		case "AUTO_INCREMENT":
			field.Extra = "auto_increment"
		case "DEFAULT":
			if i+1 < len(words) {
				defaultVal, _ := extractDefaultValue("DEFAULT " + strings.Join(words[i+1:], " "))
				field.Default = defaultVal
				if defaultVal != nil {
					if strings.HasPrefix(words[i+1], "'") {
						for j := i + 1; j < len(words); j++ {
							if strings.HasSuffix(words[j], "'") {
								i = j
								break
							}
						}
					} else {
						i++
					}
				}
			}
		case "COMMENT":
			if i+1 < len(words) {
				comment := strings.Join(words[i+1:], " ")
				comment = strings.TrimSpace(comment)
				if len(comment) >= 2 && comment[0] == '\'' && comment[len(comment)-1] == '\'' {
					comment = comment[1 : len(comment)-1]
					comment = strings.ReplaceAll(comment, "''", "'")
				}
				field.Comment = comment
				return
			}
		case "ON":
			if i+1 < len(upperWords) && upperWords[i+1] == "UPDATE" {
				if i+2 < len(upperWords) {
					if field.Extra != "" {
						field.Extra += ", "
					}
					field.Extra += "on update " + strings.ToLower(words[i+2])
					i += 2
				}
			}
		}
	}
}
