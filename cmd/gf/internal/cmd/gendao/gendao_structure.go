package gendao

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/olekukonko/tablewriter"
)

type generateStructDefinitionInput struct {
	CGenDaoInternalInput
	StructName string                     // Struct name.
	FieldMap   map[string]*gdb.TableField // Table field map.
	IsDo       bool                       // Is generating DTO struct.
}

func generateStructDefinition(in generateStructDefinitionInput) string {
	buffer := bytes.NewBuffer(nil)
	array := make([][]string, len(in.FieldMap))
	names := sortFieldKeyForDao(in.FieldMap)
	for index, name := range names {
		field := in.FieldMap[name]
		array[index] = generateStructFieldDefinition(field, in)
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	stContent := buffer.String()
	// Let's do this hack of table writer for indent!
	stContent = gstr.Replace(stContent, "  #", "")
	stContent = gstr.Replace(stContent, "` ", "`")
	stContent = gstr.Replace(stContent, "``", "")
	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", in.StructName))
	if in.IsDo {
		buffer.WriteString(fmt.Sprintf("g.Meta `orm:\"table:%s, do:true\"`\n", in.TableName))
	}
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String()
}

// generateStructFieldForModel generates and returns the attribute definition for specified field.
func generateStructFieldDefinition(field *gdb.TableField, in generateStructDefinitionInput) []string {
	var (
		typeName string
		jsonTag  = getJsonTagFromCase(field.Name, in.JsonCase)
	)
	t, _ := gregex.ReplaceString(`\(.+\)`, "", field.Type)
	t = gstr.Split(gstr.Trim(t), " ")[0]
	t = gstr.ToLower(t)

	switch t {
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		typeName = "[]byte"

	case "bit", "int", "int2", "tinyint", "small_int", "smallint", "medium_int", "mediumint", "serial":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint"
		} else {
			typeName = "int"
		}

	case "int4", "int8", "big_int", "bigint", "bigserial":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint64"
		} else {
			typeName = "int64"
		}

	// pgsql int32 slice.
	case "_int2":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "[]uint"
		} else {
			typeName = "[]int"
		}

	// pgsql int64 slice.
	case "_int4", "_int8":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "[]uint64"
		} else {
			typeName = "[]int64"
		}

	case "real":
		typeName = "float32"

	case "float", "double", "decimal", "smallmoney", "numeric":
		typeName = "float64"

	case "bool":
		typeName = "bool"

	case "datetime", "timestamp", "date", "time":
		if in.StdTime {
			typeName = "time.Time"
		} else {
			typeName = "*gtime.Time"
		}
	case "json", "jsonb":
		if in.GJsonSupport {
			typeName = "*gjson.Json"
		} else {
			typeName = "string"
		}
	default:
		// Automatically detect its data type.
		switch {
		case strings.Contains(t, "int"):
			typeName = "int"
		case strings.Contains(t, "text") || strings.Contains(t, "char"):
			typeName = "string"
		case strings.Contains(t, "float") || strings.Contains(t, "double"):
			typeName = "float64"
		case strings.Contains(t, "bool"):
			typeName = "bool"
		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			typeName = "[]byte"
		case strings.Contains(t, "date") || strings.Contains(t, "time"):
			if in.StdTime {
				typeName = "time.Time"
			} else {
				typeName = "*gtime.Time"
			}
		default:
			typeName = "string"
		}
	}

	var (
		tagKey = "`"
		result = []string{
			"    #" + gstr.CaseCamel(field.Name),
			" #" + typeName,
		}
		descriptionTag = gstr.Replace(formatComment(field.Comment), `"`, `\"`)
	)

	result = append(result, " #"+fmt.Sprintf(tagKey+`json:"%s"`, jsonTag))
	result = append(result, " #"+fmt.Sprintf(`description:"%s"`+tagKey, descriptionTag))
	result = append(result, " #"+fmt.Sprintf(`// %s`, formatComment(field.Comment)))

	for k, v := range result {
		if in.NoJsonTag {
			v, _ = gregex.ReplaceString(`json:".+"`, ``, v)
		}
		if !in.DescriptionTag {
			v, _ = gregex.ReplaceString(`description:".*"`, ``, v)
		}
		if in.NoModelComment {
			v, _ = gregex.ReplaceString(`//.+`, ``, v)
		}
		result[k] = v
	}
	return result
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

// getJsonTagFromCase call gstr.Case* function to convert the s to specified case.
func getJsonTagFromCase(str, caseStr string) string {
	switch gstr.ToLower(caseStr) {
	case gstr.ToLower("Camel"):
		return gstr.CaseCamel(str)

	case gstr.ToLower("CamelLower"):
		return gstr.CaseCamelLower(str)

	case gstr.ToLower("Kebab"):
		return gstr.CaseKebab(str)

	case gstr.ToLower("KebabScreaming"):
		return gstr.CaseKebabScreaming(str)

	case gstr.ToLower("Snake"):
		return gstr.CaseSnake(str)

	case gstr.ToLower("SnakeFirstUpper"):
		return gstr.CaseSnakeFirstUpper(str)

	case gstr.ToLower("SnakeScreaming"):
		return gstr.CaseSnakeScreaming(str)
	}
	return str
}
