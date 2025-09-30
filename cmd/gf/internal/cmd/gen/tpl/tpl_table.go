package tpl

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
)

// Table description
type Table struct {
	Name           string // 表名
	OutputName     string // 输出表名，用于生成文件名
	OutputNameCase string // 输出表名的命名规则
	PackageName    string
	db             gdb.DB
	Fields         TableFields
	FieldsSource   map[string]*gdb.TableField
	Imports        map[string]struct{}
}

type Tables []*Table

// NewTable description
//
// createTime: 2023-12-11 16:17:33
//
// author: hailaz
func NewTable(t *TplObj, tableName string) (*Table, error) {
	fields, err := t.db.TableFields(t.ctx, tableName)
	if err != nil {
		return nil, err
	}
	table := Table{
		Name:         tableName,
		OutputName:   t.TableOutputName(tableName),
		FieldsSource: fields,
		db:           t.db,
		Imports:      make(map[string]struct{}),
	}
	table.toTableFields(t.in)
	return &table, nil
}

// Name description
//
// createTime: 2023-10-23 16:17:30
//
// author: hailaz
func (t *Table) Show() string {
	return fmt.Sprintf("table name is %s", t.Name)
}

// NameCase description
func (t *Table) NameCase() string {
	return gstr.CaseConvert(t.Name, gstr.CaseTypeMatch(t.OutputNameCase))
}

// NameCaseCamel description
func (t *Table) NameCaseCamel() string {
	return gstr.CaseCamel(t.Name)
}

// NameCaseCamelLower description
func (t *Table) NameCaseCamelLower() string {
	return gstr.CaseCamelLower(t.Name)
}

// NameCaseSnake description
func (t *Table) NameCaseSnake() string {
	return gstr.CaseSnake(t.Name)
}

// NameCaseKebabScreaming description
func (t *Table) NameCaseKebabScreaming() string {
	return gstr.CaseKebabScreaming(t.Name)
}

// FileName description
func (t *Table) FileName() string {
	return gstr.CaseConvert(t.OutputName, gstr.CaseTypeMatch(t.OutputNameCase)) + ".go"
}

// toTableFields description
//
// createTime: 2023-10-23 17:22:40
//
// author: hailaz
func (t *Table) toTableFields(in CGenTplInput) {
	if len(t.Fields) > 0 {
		return
	}
	t.Fields = make(TableFields, len(t.FieldsSource))
	for _, v := range t.FieldsSource {
		field := &TableField{
			TableField: *v,
			JsonCase:   in.JsonCase,
			CustomTags: make(map[string]string),
		}

		// 设置字段类型
		appendImport := field.GetLocalTypeName(context.Background(), t.db, Input{
			TypeMapping:  in.TypeMapping,
			FieldMapping: in.FieldMapping,
			StdTime:      in.StdTime,
			GJsonSupport: in.GJsonSupport,
		})
		if appendImport != "" {
			t.Imports[appendImport] = struct{}{}
		}

		// 从 FieldMapping 中提取自定义标签
		if in.FieldMapping != nil {
			if fieldMapping, ok := in.FieldMapping[v.Name]; ok {
				if fieldMapping.Tags != nil {
					for tagName, tagValue := range fieldMapping.Tags {
						field.CustomTags[tagName] = tagValue
					}
				}
			}
		}

		t.Fields[v.Index] = field
	}
}

// SortFields 字段排序
//
// createTime: 2023-10-23 17:18:22
//
// author: hailaz
func (t *Table) SortFields(isReverse bool) {
	if isReverse {
		sort.Sort(sort.Reverse(t.Fields))
	} else {
		sort.Sort(t.Fields)
	}
}

// FieldsJsonStr 表字段json字符串
//
// createTime: 2023-10-23 17:29:39
//
// author: hailaz
func (t *Table) FieldsJsonStr(caseName string) string {
	mapStr := make(map[string]interface{}, len(t.Fields))
	for _, v := range t.Fields {
		mapStr[v.NameCaseConvert(caseName)] = v.Default
	}
	b, err := json.MarshalIndent(mapStr, "", "    ")
	if err != nil {
		return ""
	}
	return string(b)
}

// TagInput holds input for tag generation
type TagInput struct {
	in CGenTplInput
}

// GetTagInput returns TagInput for template usage
func (t *Table) GetTagInput(in CGenTplInput) TagInput {
	return TagInput{in: in}
}

// GetTables 获取数据库表结构信息
func (t *TplObj) GetTables() (Tables, error) {
	nameList, err := t.db.Tables(t.ctx)
	if err != nil {
		return nil, err
	}

	// 过滤表名
	nameList = filterTablesByName(nameList, t.in.TableNamePattern)

	// 根据Tables参数过滤
	nameList = filterTablesByInclude(nameList, t.in.Tables)

	// 根据TablesEx参数过滤
	nameList = filterTablesByExclude(nameList, t.in.TablesEx)

	tables := make(Tables, 0, len(nameList))
	for _, v := range nameList {
		t, err := NewTable(t, v)
		if err != nil {
			continue
		}

		t.SortFields(true)
		tables = append(tables, t)
	}

	return tables, nil
}

// TableOutputName description
//
// createTime: 2025-01-25 17:20:46
func (t *TplObj) TableOutputName(name string) string {
	if t.in.Prefix != "" {
		name = t.in.Prefix + name
	}

	if t.in.RemovePrefix != "" {
		name = strings.TrimPrefix(name, t.in.RemovePrefix)
	}

	return name
}

// 新增过滤函数
func filterTablesByName(tables []string, pattern string) []string {
	if pattern == "" {
		return tables
	}
	var result []string
	re, err := regexp.Compile(pattern)
	if err != nil {
		return tables
	}
	for _, table := range tables {
		if re.MatchString(table) {
			result = append(result, table)
		}
	}
	return result
}

// 根据包含表名过滤
func filterTablesByInclude(tables []string, include string) []string {
	if include == "" {
		return tables
	}
	includeTables := strings.Split(include, ",")
	result := make([]string, 0, len(includeTables))
	for _, table := range tables {
		for _, includeTable := range includeTables {
			if table == includeTable {
				result = append(result, table)
				break
			}
		}
	}
	return result
}

// 根据排除表名过滤
func filterTablesByExclude(tables []string, exclude string) []string {
	if exclude == "" {
		return tables
	}
	excludeTables := strings.Split(exclude, ",")
	result := make([]string, 0, len(tables))
	for _, table := range tables {
		exclude := false
		for _, excludeTable := range excludeTables {
			if table == excludeTable {
				exclude = true
				break
			}
		}
		if !exclude {
			result = append(result, table)
		}
	}
	return result
}
