package gendao

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/olekukonko/tablewriter"
)

type generateStructDefinitionInput struct {
	CGenDaoInternalInput
	TableName  string                     // Table name.
	StructName string                     // Struct name.
	FieldMap   map[string]*gdb.TableField // Table field map.
	IsDo       bool                       // Is generating DTO struct.
}

func generateStructDefinition(ctx context.Context, in generateStructDefinitionInput) string {
	buffer := bytes.NewBuffer(nil)
	array := make([][]string, len(in.FieldMap))
	names := sortFieldKeyForDao(in.FieldMap)
	for index, name := range names {
		field := in.FieldMap[name]
		array[index] = generateStructFieldDefinition(ctx, field, in)
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
func generateStructFieldDefinition(
	ctx context.Context, field *gdb.TableField, in generateStructDefinitionInput,
) []string {
	var (
		err      error
		typeName string
		jsonTag  = getJsonTagFromCase(field.Name, in.JsonCase)
	)
	typeName, err = in.DB.CheckLocalTypeForField(ctx, field.Type, nil)
	if err != nil {
		panic(err)
	}
	switch typeName {
	case gdb.LocalTypeDate, gdb.LocalTypeDatetime:
		if in.StdTime {
			typeName = "time.Time"
		} else {
			typeName = "*gtime.Time"
		}

	case gdb.LocalTypeInt64Bytes:
		typeName = "int64"

	case gdb.LocalTypeUint64Bytes:
		typeName = "uint64"

	// Special type handle.
	case gdb.LocalTypeJson, gdb.LocalTypeJsonb:
		if in.GJsonSupport {
			typeName = "*gjson.Json"
		} else {
			typeName = "string"
		}

	case gdb.LocalTypeMap:
		typeName = "g.Map"

	case gdb.LocalTypeArray:
		typeName = "garray.Array"

	case gdb.LocalTypeUUID:
		typeName = "uuid.UUID"

	case gdb.LocalTypeDecimal:
		typeName = "decimal.Decimal"

	case gdb.LocalTypeBigInt:
		typeName = "*big.Int"

	case gdb.LocalTypeInterface:
		typeName = "interface{}"

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
