// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/v2/internal/structs"
	"github.com/gogf/gf/v2/text/gstr"
	"reflect"
	"strings"
)

// iString is used for type assert api for String().
type iString interface {
	String() string
}

// iMarshalJSON is the interface for custom Json marshaling.
type iMarshalJSON interface {
	MarshalJSON() ([]byte, error)
}

// ExportOption specifies the behavior of function Export.
type ExportOption struct {
	WithoutType bool // WithoutType specifies exported content has no type information.
}

// Dump prints variables `values` to stdout with more manually readable.
func Dump(values ...interface{}) {
	for _, value := range values {
		if s := Export(value, ExportOption{
			WithoutType: true,
		}); s != "" {
			fmt.Println(s)
		}
	}
}

// DumpWithType acts like Dump, but with type information.
// Also see Dump.
func DumpWithType(values ...interface{}) {
	for _, value := range values {
		if s := Export(value, ExportOption{
			WithoutType: false,
		}); s != "" {
			fmt.Println(s)
		}
	}
}

// Export returns variables `values` as a string with more manually readable.
func Export(value interface{}, option ExportOption) string {
	buffer := bytes.NewBuffer(nil)
	doExport(value, "", buffer, doExportOption{
		WithoutType: option.WithoutType,
	})
	return buffer.String()
}

type doExportOption struct {
	WithoutType bool
}

func doExport(value interface{}, indent string, buffer *bytes.Buffer, option doExportOption) {
	if value == nil {
		buffer.WriteString(`<nil>`)
		return
	}
	var (
		reflectValue    = reflect.ValueOf(value)
		reflectKind     = reflectValue.Kind()
		reflectTypeName = reflect.TypeOf(value).String()
		newIndent       = indent + dumpIndent
	)
	reflectTypeName = strings.ReplaceAll(reflectTypeName, `[]uint8`, `[]byte`)
	if option.WithoutType {
		reflectTypeName = ""
	}
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		if _, ok := value.([]byte); ok {
			if option.WithoutType {
				buffer.WriteString(fmt.Sprintf(`"%s"`, value))
			} else {
				buffer.WriteString(fmt.Sprintf(
					`%s(%d) "%s"`,
					reflectTypeName,
					len(reflectValue.String()),
					value,
				))
			}
			return
		}
		if reflectValue.Len() == 0 {
			if option.WithoutType {
				buffer.WriteString("[]")
			} else {
				buffer.WriteString(fmt.Sprintf("%s(0) []", reflectTypeName))
			}
			return
		}
		if option.WithoutType {
			buffer.WriteString("[\n")
		} else {
			buffer.WriteString(fmt.Sprintf("%s(%d) [\n", reflectTypeName, reflectValue.Len()))
		}
		for i := 0; i < reflectValue.Len(); i++ {
			buffer.WriteString(newIndent)
			doExport(reflectValue.Index(i).Interface(), newIndent, buffer, option)
			buffer.WriteString(",\n")
		}
		buffer.WriteString(fmt.Sprintf("%s]", indent))

	case reflect.Map:
		var (
			mapKeys = reflectValue.MapKeys()
		)
		if len(mapKeys) == 0 {
			if option.WithoutType {
				buffer.WriteString("{}")
			} else {
				buffer.WriteString(fmt.Sprintf("%s(0) {}", reflectTypeName))
			}
			return
		}

		var (
			maxSpaceNum = 0
			tmpSpaceNum = 0
			mapKeyStr   = ""
		)
		for _, key := range mapKeys {
			tmpSpaceNum = len(fmt.Sprintf(`%v`, key.Interface()))
			if tmpSpaceNum > maxSpaceNum {
				maxSpaceNum = tmpSpaceNum
			}
		}
		if option.WithoutType {
			buffer.WriteString("{\n")
		} else {
			buffer.WriteString(fmt.Sprintf("%s(%d) {\n", reflectTypeName, len(mapKeys)))
		}
		for _, mapKey := range mapKeys {
			tmpSpaceNum = len(fmt.Sprintf(`%v`, mapKey.Interface()))
			if mapKey.Kind() == reflect.String {
				mapKeyStr = fmt.Sprintf(`"%v"`, mapKey.Interface())
			} else {
				mapKeyStr = fmt.Sprintf(`%v`, mapKey.Interface())
			}
			if option.WithoutType {
				buffer.WriteString(fmt.Sprintf(
					"%s%v:%s",
					newIndent,
					mapKeyStr,
					strings.Repeat(" ", maxSpaceNum-tmpSpaceNum+1),
				))
			} else {
				buffer.WriteString(fmt.Sprintf(
					"%s%s(%v):%s",
					newIndent,
					mapKey.Type().String(),
					mapKeyStr,
					strings.Repeat(" ", maxSpaceNum-tmpSpaceNum+1),
				))
			}
			doExport(reflectValue.MapIndex(mapKey).Interface(), newIndent, buffer, option)
			buffer.WriteString(",\n")
		}
		buffer.WriteString(fmt.Sprintf("%s}", indent))

	case reflect.Struct:
		structFields, _ := structs.Fields(structs.FieldsInput{
			Pointer:         value,
			RecursiveOption: structs.RecursiveOptionEmbeddedNoTag,
		})
		if len(structFields) == 0 {
			var (
				structContentStr  = ""
				attributeCountStr = "0"
			)
			if v, ok := value.(iString); ok {
				structContentStr = v.String()
			} else if v, ok := value.(iMarshalJSON); ok {
				b, _ := v.MarshalJSON()
				structContentStr = string(b)
			}
			if structContentStr == "" {
				structContentStr = "{}"
			} else {
				if strings.HasPrefix(structContentStr, `"`) && strings.HasSuffix(structContentStr, `"`) {
					attributeCountStr = fmt.Sprintf(`%d`, len(structContentStr))
				} else {
					structContentStr = fmt.Sprintf(`"%s"`, gstr.AddSlashes(structContentStr))
					attributeCountStr = fmt.Sprintf(`%d`, len(structContentStr)-2)
				}
			}
			if option.WithoutType {
				buffer.WriteString(structContentStr)
			} else {
				buffer.WriteString(fmt.Sprintf(
					"%s(%s) %s",
					reflectTypeName,
					attributeCountStr,
					structContentStr,
				))
			}
			return
		}

		var (
			maxSpaceNum = 0
			tmpSpaceNum = 0
		)
		for _, field := range structFields {
			tmpSpaceNum = len(field.Name())
			if tmpSpaceNum > maxSpaceNum {
				maxSpaceNum = tmpSpaceNum
			}
		}
		if option.WithoutType {
			buffer.WriteString("{\n")
		} else {
			buffer.WriteString(fmt.Sprintf("%s(%d) {\n", reflectTypeName, len(structFields)))
		}
		for _, field := range structFields {
			tmpSpaceNum = len(fmt.Sprintf(`%v`, field.Name()))
			buffer.WriteString(fmt.Sprintf(
				"%s%s:%s",
				newIndent,
				field.Name(),
				strings.Repeat(" ", maxSpaceNum-tmpSpaceNum+1),
			))
			doExport(field.Value.Interface(), newIndent, buffer, option)
			buffer.WriteString(",\n")
		}
		buffer.WriteString(fmt.Sprintf("%s}", indent))

	case reflect.String:
		if option.WithoutType {
			buffer.WriteString(fmt.Sprintf("\"%v\"", value))
		} else {
			buffer.WriteString(fmt.Sprintf(
				"%s(%d) \"%v\"",
				reflectTypeName,
				len(reflectValue.String()),
				value,
			))
		}

	default:
		if option.WithoutType {
			buffer.WriteString(fmt.Sprintf("%v", value))
		} else {
			buffer.WriteString(fmt.Sprintf("%s(%v)", reflectTypeName, value))
		}
	}
}
