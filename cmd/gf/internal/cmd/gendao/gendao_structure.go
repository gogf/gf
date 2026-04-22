// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// generateStructDefinitionInput holds parameters for generating a Go struct definition
// from database table fields.
type generateStructDefinitionInput struct {
	CGenDaoInternalInput
	TableName  string                     // Original database table name.
	StructName string                     // Go struct name (CamelCase of table name).
	FieldMap   map[string]*gdb.TableField // Map of column name to field metadata.
	IsDo       bool                       // Whether generating a DO struct (uses g.Meta orm tag).
}

// generateStructDefinition generates a complete Go struct definition string from table fields.
// It returns the struct source code and a list of additional import paths needed
// by custom type mappings. The fields are rendered in a table-aligned format
// using tablewriter for consistent code formatting.
func generateStructDefinition(ctx context.Context, in generateStructDefinitionInput) (string, []string) {
	var appendImports []string
	buffer := bytes.NewBuffer(nil)
	array := make([][]string, len(in.FieldMap))
	names := sortFieldKeyForDao(in.FieldMap)
	for index, name := range names {
		var imports string
		field := in.FieldMap[name]
		array[index], imports = generateStructFieldDefinition(ctx, field, in)
		if imports != "" {
			appendImports = append(appendImports, imports)
		}
	}
	table := tablewriter.NewTable(buffer, twRenderer, twConfig)
	table.Bulk(array)
	table.Render()
	stContent := buffer.String()
	// Let's do this hack of table writer for indent!
	stContent = gstr.Replace(stContent, "  #", "")
	stContent = gstr.Replace(stContent, "` ", "`")
	stContent = gstr.Replace(stContent, "``", "")
	buffer.Reset()
	fmt.Fprintf(buffer, "type %s struct {\n", in.StructName)
	if in.IsDo {
		fmt.Fprintf(buffer, "g.Meta `orm:\"table:%s, do:true\"`\n", in.TableName)
	}
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String(), appendImports
}

// getTypeMappingInfo looks up a database field type in the type mapping configuration.
// It handles exact matches first, then tries to extract the base type name from
// parameterized types like "varchar(255)" or "numeric(10,2) unsigned".
// Returns the mapped Go type name and its import path (if any).
func getTypeMappingInfo(
	ctx context.Context, fieldType string, inTypeMapping map[DBFieldTypeName]CustomAttributeType,
) (typeNameStr, importStr string) {
	if typeMapping, ok := inTypeMapping[strings.ToLower(fieldType)]; ok {
		typeNameStr = typeMapping.Type
		importStr = typeMapping.Import
		return
	}
	tryTypeMatch, _ := gregex.MatchString(`(.+?)\(([^\(\)]+)\)([\s\)]*)`, fieldType)
	var (
		tryTypeName string
		moreTry     bool
	)
	if len(tryTypeMatch) == 4 {
		tryTypeMatch3, _ := gregex.ReplaceString(`\s+`, "", tryTypeMatch[3])
		tryTypeName = gstr.Trim(tryTypeMatch[1]) + tryTypeMatch3
		moreTry = tryTypeMatch3 != ""
	} else {
		tryTypeName = gstr.Split(fieldType, " ")[0]
	}
	if tryTypeName != "" {
		if typeMapping, ok := inTypeMapping[strings.ToLower(tryTypeName)]; ok {
			typeNameStr = typeMapping.Type
			importStr = typeMapping.Import
		} else if moreTry {
			typeNameStr, importStr = getTypeMappingInfo(ctx, tryTypeName, inTypeMapping)
		}
	}
	return
}

// generateStructFieldDefinition generates and returns the attribute definition for specified field.
func generateStructFieldDefinition(
	ctx context.Context, field *gdb.TableField, in generateStructDefinitionInput,
) (attrLines []string, appendImport string) {
	var (
		err              error
		localTypeName    gdb.LocalType
		localTypeNameStr string
	)

	if in.TypeMapping != nil && len(in.TypeMapping) > 0 {
		localTypeNameStr, appendImport = getTypeMappingInfo(ctx, field.Type, in.TypeMapping)
	}

	if localTypeNameStr == "" {
		if in.DB != nil {
			localTypeName, err = in.DB.CheckLocalTypeForField(ctx, field.Type, nil)
			if err != nil {
				panic(err)
			}
		} else {
			// SQL file mode: use standalone type checking without database connection.
			localTypeName, err = gdb.CheckLocalTypeForFieldType(field.Type)
			if err != nil {
				panic(err)
			}
		}
		localTypeNameStr = string(localTypeName)
		switch localTypeName {
		case gdb.LocalTypeDate, gdb.LocalTypeTime, gdb.LocalTypeDatetime:
			if in.StdTime {
				localTypeNameStr = "time.Time"
			} else {
				localTypeNameStr = "*gtime.Time"
			}

		case gdb.LocalTypeInt64Bytes:
			localTypeNameStr = "int64"

		case gdb.LocalTypeUint64Bytes:
			localTypeNameStr = "uint64"

		// Special type handle.
		case gdb.LocalTypeJson, gdb.LocalTypeJsonb:
			if in.GJsonSupport {
				localTypeNameStr = "*gjson.Json"
			} else {
				localTypeNameStr = "string"
			}
		}
	}

	var (
		tagKey         = "`"
		descriptionTag = gstr.Replace(formatComment(field.Comment), `"`, `\"`)
	)
	removeFieldPrefixArray := gstr.SplitAndTrim(in.RemoveFieldPrefix, ",")
	newFiledName := field.Name
	for _, v := range removeFieldPrefixArray {
		newFiledName = gstr.TrimLeftStr(newFiledName, v, 1)
	}

	if in.FieldMapping != nil && len(in.FieldMapping) > 0 {
		if typeMapping, ok := in.FieldMapping[fmt.Sprintf("%s.%s", in.TableName, newFiledName)]; ok {
			localTypeNameStr = typeMapping.Type
			appendImport = typeMapping.Import
		}
	}

	attrLines = []string{
		"    #" + formatFieldName(newFiledName, FieldNameCaseCamel),
		" #" + localTypeNameStr,
	}

	jsonTag := gstr.CaseConvert(newFiledName, gstr.CaseTypeMatch(in.JsonCase))
	attrLines = append(attrLines, fmt.Sprintf(` #%sjson:"%s"`, tagKey, jsonTag))
	// orm tag
	if !in.IsDo {
		// entity
		attrLines = append(attrLines, fmt.Sprintf(` #orm:"%s"`, field.Name))
	}
	attrLines = append(attrLines, fmt.Sprintf(` #description:"%s"%s`, descriptionTag, tagKey))
	attrLines = append(attrLines, fmt.Sprintf(` #// %s`, formatComment(field.Comment)))

	for k, v := range attrLines {
		if in.NoJsonTag {
			v, _ = gregex.ReplaceString(`json:".+"`, ``, v)
		}
		if !in.DescriptionTag {
			v, _ = gregex.ReplaceString(`description:".*"`, ``, v)
		}
		if in.NoModelComment {
			v, _ = gregex.ReplaceString(`//.+`, ``, v)
		}
		attrLines[k] = v
	}
	return attrLines, appendImport
}

// FieldNameCase defines the naming convention for converting field names to Go identifiers.
type FieldNameCase string

const (
	FieldNameCaseCamel      FieldNameCase = "CaseCamel"      // PascalCase: "user_name" -> "UserName"
	FieldNameCaseCamelLower FieldNameCase = "CaseCamelLower" // camelCase: "user_name" -> "userName"
)

// formatFieldName formats and returns a new field name that is used for golang codes generating.
func formatFieldName(fieldName string, nameCase FieldNameCase) string {
	// For normal databases like mysql, pgsql, sqlite,
	// field/table names of that are in normal case.
	var newFieldName = fieldName
	if isAllUpper(fieldName) {
		// For special databases like dm, oracle,
		// field/table names of that are in upper case.
		newFieldName = strings.ToLower(fieldName)
	}
	switch nameCase {
	case FieldNameCaseCamel:
		return gstr.CaseCamel(newFieldName)
	case FieldNameCaseCamelLower:
		return gstr.CaseCamelLower(newFieldName)
	default:
		return ""
	}
}

// isAllUpper checks and returns whether given `fieldName` all letters are upper case.
func isAllUpper(fieldName string) bool {
	for _, b := range fieldName {
		if b >= 'a' && b <= 'z' {
			return false
		}
	}
	return true
}

// formatComment formats the comment string to fit the golang code without any lines.
func formatComment(comment string) string {
	comment = gstr.ReplaceByArray(comment, g.SliceStr{
		"\n", " ",
		"\r", " ",
	})
	comment = gstr.Replace(comment, `\n`, " ")
	comment = gstr.Trim(comment)
	return comment
}
