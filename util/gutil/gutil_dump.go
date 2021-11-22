// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/internal/structs"
	"github.com/gogf/gf/v2/text/gstr"
)

// iString is used for type assert api for String().
type iString interface {
	String() string
}

// iError is used for type assert api for Error().
type iError interface {
	Error() string
}

// iMarshalJSON is the interface for custom Json marshaling.
type iMarshalJSON interface {
	MarshalJSON() ([]byte, error)
}

// DumpOption specifies the behavior of function Export.
type DumpOption struct {
	WithType bool // WithType specifies dumping content with type information.
}

// Dump prints variables `values` to stdout with more manually readable.
func Dump(values ...interface{}) {
	for _, value := range values {
		DumpWithOption(value, DumpOption{
			WithType: false,
		})
	}
}

// DumpWithType acts like Dump, but with type information.
// Also see Dump.
func DumpWithType(values ...interface{}) {
	for _, value := range values {
		DumpWithOption(value, DumpOption{
			WithType: true,
		})
	}
}

// DumpWithOption returns variables `values` as a string with more manually readable.
func DumpWithOption(value interface{}, option DumpOption) {
	buffer := bytes.NewBuffer(nil)
	DumpTo(buffer, value, DumpOption{
		WithType: option.WithType,
	})
	fmt.Println(buffer.String())
}

// DumpTo writes variables `values` as a string in to `writer` with more manually readable
func DumpTo(writer io.Writer, value interface{}, option DumpOption) {
	buffer := bytes.NewBuffer(nil)
	doDump(value, "", buffer, doDumpOption{
		WithType: option.WithType,
	})
	_, _ = writer.Write(buffer.Bytes())
}

type doDumpOption struct {
	WithType bool
}

func doDump(value interface{}, indent string, buffer *bytes.Buffer, option doDumpOption) {
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
	if !option.WithType {
		reflectTypeName = ""
	}
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	var (
		exportInternalInput = doDumpInternalInput{
			Value:           value,
			Indent:          indent,
			NewIndent:       newIndent,
			Buffer:          buffer,
			Option:          option,
			ReflectValue:    reflectValue,
			ReflectTypeName: reflectTypeName,
		}
	)
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		doDumpSlice(exportInternalInput)

	case reflect.Map:
		doDumpMap(exportInternalInput)

	case reflect.Struct:
		doDumpStruct(exportInternalInput)

	case reflect.String:
		doDumpString(exportInternalInput)

	case reflect.Bool:
		if reflectValue.Bool() {
			buffer.WriteString(`true`)
		} else {
			buffer.WriteString(`false`)
		}

	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		doDumpNumber(exportInternalInput)

	case reflect.Chan:
		buffer.WriteString(`<chan>`)

	case reflect.Func:
		buffer.WriteString(`<func>`)

	default:
		doDumpDefault(exportInternalInput)
	}
}

type doDumpInternalInput struct {
	Value           interface{}
	Indent          string
	NewIndent       string
	Buffer          *bytes.Buffer
	Option          doDumpOption
	ReflectValue    reflect.Value
	ReflectTypeName string
}

func doDumpSlice(in doDumpInternalInput) {
	if b, ok := in.Value.([]byte); ok {
		if !in.Option.WithType {
			in.Buffer.WriteString(fmt.Sprintf(`"%s"`, addSlashesForString(string(b))))
		} else {
			in.Buffer.WriteString(fmt.Sprintf(
				`%s(%d) "%s"`,
				in.ReflectTypeName,
				len(string(b)),
				string(b),
			))
		}
		return
	}
	if in.ReflectValue.Len() == 0 {
		if !in.Option.WithType {
			in.Buffer.WriteString("[]")
		} else {
			in.Buffer.WriteString(fmt.Sprintf("%s(0) []", in.ReflectTypeName))
		}
		return
	}
	if !in.Option.WithType {
		in.Buffer.WriteString("[\n")
	} else {
		in.Buffer.WriteString(fmt.Sprintf("%s(%d) [\n", in.ReflectTypeName, in.ReflectValue.Len()))
	}
	for i := 0; i < in.ReflectValue.Len(); i++ {
		in.Buffer.WriteString(in.NewIndent)
		doDump(in.ReflectValue.Index(i).Interface(), in.NewIndent, in.Buffer, in.Option)
		in.Buffer.WriteString(",\n")
	}
	in.Buffer.WriteString(fmt.Sprintf("%s]", in.Indent))
}

func doDumpMap(in doDumpInternalInput) {
	var (
		mapKeys = in.ReflectValue.MapKeys()
	)
	if len(mapKeys) == 0 {
		if !in.Option.WithType {
			in.Buffer.WriteString("{}")
		} else {
			in.Buffer.WriteString(fmt.Sprintf("%s(0) {}", in.ReflectTypeName))
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
	if !in.Option.WithType {
		in.Buffer.WriteString("{\n")
	} else {
		in.Buffer.WriteString(fmt.Sprintf("%s(%d) {\n", in.ReflectTypeName, len(mapKeys)))
	}
	for _, mapKey := range mapKeys {
		tmpSpaceNum = len(fmt.Sprintf(`%v`, mapKey.Interface()))
		if mapKey.Kind() == reflect.String {
			mapKeyStr = fmt.Sprintf(`"%v"`, mapKey.Interface())
		} else {
			mapKeyStr = fmt.Sprintf(`%v`, mapKey.Interface())
		}
		// Map key and indent string dump.
		if !in.Option.WithType {
			in.Buffer.WriteString(fmt.Sprintf(
				"%s%v:%s",
				in.NewIndent,
				mapKeyStr,
				strings.Repeat(" ", maxSpaceNum-tmpSpaceNum+1),
			))
		} else {
			in.Buffer.WriteString(fmt.Sprintf(
				"%s%s(%v):%s",
				in.NewIndent,
				mapKey.Type().String(),
				mapKeyStr,
				strings.Repeat(" ", maxSpaceNum-tmpSpaceNum+1),
			))
		}
		// Map value dump.
		doDump(in.ReflectValue.MapIndex(mapKey).Interface(), in.NewIndent, in.Buffer, in.Option)
		in.Buffer.WriteString(",\n")
	}
	in.Buffer.WriteString(fmt.Sprintf("%s}", in.Indent))
}

func doDumpStruct(in doDumpInternalInput) {
	structFields, _ := structs.Fields(structs.FieldsInput{
		Pointer:         in.Value,
		RecursiveOption: structs.RecursiveOptionEmbedded,
	})
	if len(structFields) == 0 {
		var (
			structContentStr  = ""
			attributeCountStr = "0"
		)
		if v, ok := in.Value.(iString); ok {
			structContentStr = v.String()
		} else if v, ok := in.Value.(iError); ok {
			structContentStr = v.Error()
		} else if v, ok := in.Value.(iMarshalJSON); ok {
			b, _ := v.MarshalJSON()
			structContentStr = string(b)
		}
		if structContentStr == "" {
			structContentStr = "{}"
		} else {
			structContentStr = fmt.Sprintf(`"%s"`, addSlashesForString(structContentStr))
			attributeCountStr = fmt.Sprintf(`%d`, len(structContentStr)-2)
		}
		if !in.Option.WithType {
			in.Buffer.WriteString(structContentStr)
		} else {
			in.Buffer.WriteString(fmt.Sprintf(
				"%s(%s) %s",
				in.ReflectTypeName,
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
	if !in.Option.WithType {
		in.Buffer.WriteString("{\n")
	} else {
		in.Buffer.WriteString(fmt.Sprintf("%s(%d) {\n", in.ReflectTypeName, len(structFields)))
	}
	for _, field := range structFields {
		tmpSpaceNum = len(fmt.Sprintf(`%v`, field.Name()))
		in.Buffer.WriteString(fmt.Sprintf(
			"%s%s:%s",
			in.NewIndent,
			field.Name(),
			strings.Repeat(" ", maxSpaceNum-tmpSpaceNum+1),
		))
		doDump(field.Value.Interface(), in.NewIndent, in.Buffer, in.Option)
		in.Buffer.WriteString(",\n")
	}
	in.Buffer.WriteString(fmt.Sprintf("%s}", in.Indent))
}

func doDumpNumber(in doDumpInternalInput) {
	if v, ok := in.Value.(iString); ok {
		s := v.String()
		if !in.Option.WithType {
			in.Buffer.WriteString(fmt.Sprintf(`"%v"`, addSlashesForString(s)))
		} else {
			in.Buffer.WriteString(fmt.Sprintf(
				`%s(%d) "%v"`,
				in.ReflectTypeName,
				len(s),
				addSlashesForString(s),
			))
		}
	} else {
		doDumpDefault(in)
	}
}

func doDumpString(in doDumpInternalInput) {
	s := in.ReflectValue.String()
	if !in.Option.WithType {
		in.Buffer.WriteString(fmt.Sprintf(`"%v"`, addSlashesForString(s)))
	} else {
		in.Buffer.WriteString(fmt.Sprintf(
			`%s(%d) "%v"`,
			in.ReflectTypeName,
			len(s),
			addSlashesForString(s),
		))
	}
}

func doDumpDefault(in doDumpInternalInput) {
	s := fmt.Sprintf("%v", in.Value)
	s = gstr.Trim(s, `<>`)
	if !in.Option.WithType {
		in.Buffer.WriteString(s)
	} else {
		in.Buffer.WriteString(fmt.Sprintf("%s(%s)", in.ReflectTypeName, s))
	}
}

func addSlashesForString(s string) string {
	return gstr.ReplaceByMap(s, map[string]string{
		`"`:  `\"`,
		"\r": `\r`,
		"\t": `\t`,
		"\n": `\n`,
	})
}
