package tpl

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// TableField description
type TableField struct {
	gdb.TableField
	LocalType  string
	JsonCase   string
	CustomTags map[string]string // 自定义标签
}

type TableFields []*TableField

// Len returns the length of TableFields slice
func (s TableFields) Len() int { return len(s) }

// Swap swaps the elements with indexes i and j
func (s TableFields) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less reports whether the element with index i should sort before the element with index j
func (s TableFields) Less(i, j int) bool {
	return strings.Compare(s[i].Name, s[j].Name) < 0
}

// Input description
type Input struct {
	StdTime      bool
	GJsonSupport bool
	TypeMapping  map[string]CustomAttributeType
	FieldMapping map[string]CustomAttributeType
}

// GetLocalTypeName description
//
// createTime: 2023-10-25 15:43:06
//
// author: hailaz
func (field *TableField) GetLocalTypeName(ctx context.Context, db gdb.DB, in Input) (appendImport string) {
	var (
		err              error
		localTypeName    gdb.LocalType
		localTypeNameStr string
	)
	if in.TypeMapping != nil && len(in.TypeMapping) > 0 {
		var (
			tryTypeName string
		)
		tryTypeMatch, _ := gregex.MatchString(`(.+?)\((.+)\)`, field.Type)
		if len(tryTypeMatch) == 3 {
			tryTypeName = gstr.Trim(tryTypeMatch[1])
		} else {
			tryTypeName = gstr.Split(field.Type, " ")[0]
		}
		if tryTypeName != "" {
			if typeMapping, ok := in.TypeMapping[strings.ToLower(tryTypeName)]; ok {
				localTypeNameStr = typeMapping.Type
				appendImport = typeMapping.Import
			}
		}
	}

	if localTypeNameStr == "" {
		localTypeName, err = db.CheckLocalTypeForField(ctx, field.Type, nil)
		if err != nil {
			panic(err)
		}
		localTypeNameStr = string(localTypeName)
		switch localTypeName {
		case gdb.LocalTypeDate, gdb.LocalTypeDatetime:
			if in.StdTime {
				localTypeNameStr = "time.Time"
			} else {
				localTypeNameStr = "*gtime.Time"
				appendImport = "github.com/gogf/gf/v2/os/gtime"
			}

		case gdb.LocalTypeInt64Bytes:
			localTypeNameStr = "int64"

		case gdb.LocalTypeUint64Bytes:
			localTypeNameStr = "uint64"

		// Special type handle.
		case gdb.LocalTypeJson, gdb.LocalTypeJsonb:
			if in.GJsonSupport {
				localTypeNameStr = "*gjson.Json"
				appendImport = "github.com/gogf/gf/v2/encoding/gjson"
			} else {
				localTypeNameStr = "string"
			}
		}
	}

	// Check field-specific mapping (overrides type mapping)
	if len(in.FieldMapping) > 0 {
		fieldKey := field.Name
		if typeMapping, ok := in.FieldMapping[fieldKey]; ok {
			localTypeNameStr = typeMapping.Type
			if typeMapping.Import != "" {
				appendImport = typeMapping.Import
			}
		}
	}

	field.LocalType = localTypeNameStr
	return
}

// NameJsonCase description
//
// createTime: 2025-01-25 15:27:01
func (f *TableField) NameJsonCase() string {
	return gstr.CaseConvert(f.Name, gstr.CaseTypeMatch(f.JsonCase))
}

// NameCaseConvert 字段名转换
func (f *TableField) NameCaseConvert(caseName string) string {
	return gstr.CaseConvert(f.Name, gstr.CaseTypeMatch(caseName))
}

// NameCaseCamel returns the field name in camel case format
func (f *TableField) NameCaseCamel() string {
	return gstr.CaseCamel(f.Name)
}

// NameCaseCamelLower returns the field name in lower camel case format
func (f *TableField) NameCaseCamelLower() string {
	return gstr.CaseCamelLower(f.Name)
}

// NameCaseSnake returns the field name in snake case format
func (f *TableField) NameCaseSnake() string {
	return gstr.CaseSnake(f.Name)
}

// NameCaseKebabScreaming returns the field name in screaming kebab case format
func (f *TableField) NameCaseKebabScreaming() string {
	return gstr.CaseKebabScreaming(f.Name)
}

// IsNullable returns whether the field is nullable
func (f *TableField) IsNullable() bool {
	return f.Null
}

// JsonTag generates json tag for the field
func (f *TableField) JsonTag(omitempty bool, omitemptyAuto bool) string {
	if f.CustomTags != nil {
		if jsonTag, ok := f.CustomTags["json"]; ok {
			return jsonTag
		}
	}

	name := f.NameJsonCase()
	if omitempty || (omitemptyAuto && f.IsNullable()) {
		return name + ",omitempty"
	}
	return name
}

// OrmTag generates orm tag for the field
func (f *TableField) OrmTag() string {
	if f.CustomTags != nil {
		if ormTag, ok := f.CustomTags["orm"]; ok {
			return ormTag
		}
	}
	return f.Name
}

// DescriptionTag generates description tag for the field
func (f *TableField) DescriptionTag() string {
	if f.CustomTags != nil {
		if descTag, ok := f.CustomTags["description"]; ok {
			return descTag
		}
	}
	// 转义双引号
	comment := strings.ReplaceAll(f.Comment, `"`, `\"`)
	return comment
}

// CustomTag returns custom tag value by name
func (f *TableField) CustomTag(name string) string {
	if f.CustomTags == nil {
		return ""
	}
	return f.CustomTags[name]
}

// BuildTags builds all tags for the field
func (f *TableField) BuildTags(in TagBuildInput) string {
	var tags []string

	// JSON tag
	if !in.NoJsonTag {
		jsonValue := f.JsonTag(in.JsonOmitempty, in.JsonOmitemptyAuto)
		tags = append(tags, fmt.Sprintf(`json:"%s"`, jsonValue))
	}

	// ORM tag
	if in.WithOrmTag {
		ormValue := f.OrmTag()
		tags = append(tags, fmt.Sprintf(`orm:"%s"`, ormValue))
	}

	// Description tag
	if in.DescriptionTag {
		descValue := f.DescriptionTag()
		tags = append(tags, fmt.Sprintf(`description:"%s"`, descValue))
	}

	// Custom tags from CustomTags map
	if f.CustomTags != nil {
		// 按字母顺序遍历,确保输出稳定
		var keys []string
		for k := range f.CustomTags {
			// 跳过已处理的标准标签
			if k == "json" || k == "orm" || k == "description" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := f.CustomTags[k]
			tags = append(tags, fmt.Sprintf(`%s:"%s"`, k, v))
		}
	}

	if len(tags) == 0 {
		return ""
	}

	return "`" + strings.Join(tags, " ") + "`"
}

// TagBuildInput for building tags
type TagBuildInput struct {
	NoJsonTag        bool
	JsonOmitempty    bool
	JsonOmitemptyAuto bool
	WithOrmTag       bool
	DescriptionTag   bool
}
