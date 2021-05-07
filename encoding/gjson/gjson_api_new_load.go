// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/gogf/gf/internal/json"

	"github.com/gogf/gf/encoding/gini"
	"github.com/gogf/gf/encoding/gtoml"
	"github.com/gogf/gf/encoding/gxml"
	"github.com/gogf/gf/encoding/gyaml"
	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// New creates a Json object with any variable type of <data>, but <data> should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter <safe> specifies whether using this Json object in concurrent-safe context,
// which is false in default.
func New(data interface{}, safe ...bool) *Json {
	return NewWithTag(data, "json", safe...)
}

// NewWithTag creates a Json object with any variable type of <data>, but <data> should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter <tags> specifies priority tags for struct conversion to map, multiple tags joined
// with char ','.
//
// The parameter <safe> specifies whether using this Json object in concurrent-safe context, which
// is false in default.
func NewWithTag(data interface{}, tags string, safe ...bool) *Json {
	option := Option{
		Tags: tags,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return NewWithOption(data, option)
}

// NewWithOption creates a Json object with any variable type of <data>, but <data> should be a map
// or slice for data access reason, or it will make no sense.
func NewWithOption(data interface{}, option Option) *Json {
	var j *Json
	switch data.(type) {
	case string, []byte:
		if r, err := loadContentWithOption(data, option); err == nil {
			j = r
		} else {
			j = &Json{
				p:  &data,
				c:  byte(defaultSplitChar),
				vc: false,
			}
		}
	default:
		var (
			rv   = reflect.ValueOf(data)
			kind = rv.Kind()
		)
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Slice, reflect.Array:
			i := interface{}(nil)
			i = gconv.Interfaces(data)
			j = &Json{
				p:  &i,
				c:  byte(defaultSplitChar),
				vc: false,
			}
		case reflect.Map, reflect.Struct:
			i := interface{}(nil)
			i = gconv.MapDeep(data, option.Tags)
			j = &Json{
				p:  &i,
				c:  byte(defaultSplitChar),
				vc: false,
			}
		default:
			j = &Json{
				p:  &data,
				c:  byte(defaultSplitChar),
				vc: false,
			}
		}
	}
	j.mu = rwmutex.New(option.Safe)
	return j
}

// Load loads content from specified file <path>, and creates a Json object from its content.
func Load(path string, safe ...bool) (*Json, error) {
	if p, err := gfile.Search(path); err != nil {
		return nil, err
	} else {
		path = p
	}
	option := Option{}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOption(gfile.Ext(path), gfile.GetBytesWithCache(path), option)
}

// LoadJson creates a Json object from given JSON format content.
func LoadJson(data interface{}, safe ...bool) (*Json, error) {
	option := Option{}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOption("json", gconv.Bytes(data), option)
}

// LoadXml creates a Json object from given XML format content.
func LoadXml(data interface{}, safe ...bool) (*Json, error) {
	option := Option{}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOption("xml", gconv.Bytes(data), option)
}

// LoadIni creates a Json object from given INI format content.
func LoadIni(data interface{}, safe ...bool) (*Json, error) {
	option := Option{}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOption("ini", gconv.Bytes(data), option)
}

// LoadYaml creates a Json object from given YAML format content.
func LoadYaml(data interface{}, safe ...bool) (*Json, error) {
	option := Option{}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOption("yaml", gconv.Bytes(data), option)
}

// LoadToml creates a Json object from given TOML format content.
func LoadToml(data interface{}, safe ...bool) (*Json, error) {
	option := Option{}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOption("toml", gconv.Bytes(data), option)
}

// LoadContent creates a Json object from given content, it checks the data type of <content>
// automatically, supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContent(data interface{}, safe ...bool) (*Json, error) {
	content := gconv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}
	return LoadContentType(checkDataType(content), content, safe...)
}

// LoadContentType creates a Json object from given type and content,
// supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContentType(dataType string, data interface{}, safe ...bool) (*Json, error) {
	content := gconv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}
	//ignore UTF8-BOM
	if content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}
	option := Option{}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOption(dataType, content, option)
}

// IsValidDataType checks and returns whether given <dataType> a valid data type for loading.
func IsValidDataType(dataType string) bool {
	if dataType == "" {
		return false
	}
	if dataType[0] == '.' {
		dataType = dataType[1:]
	}
	switch dataType {
	case "json", "js", "xml", "yaml", "yml", "toml", "ini":
		return true
	}
	return false
}

func loadContentWithOption(data interface{}, option Option) (*Json, error) {
	content := gconv.Bytes(data)
	if len(content) == 0 {
		return NewWithOption(nil, option), nil
	}
	return loadContentTypeWithOption(checkDataType(content), content, option)
}

func loadContentTypeWithOption(dataType string, data interface{}, option Option) (*Json, error) {
	content := gconv.Bytes(data)
	if len(content) == 0 {
		return NewWithOption(nil, option), nil
	}
	//ignore UTF8-BOM
	if content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}
	return doLoadContentWithOption(dataType, content, option)
}

// doLoadContent creates a Json object from given content.
// It supports data content type as follows:
// JSON, XML, INI, YAML and TOML.
func doLoadContentWithOption(dataType string, data []byte, option Option) (*Json, error) {
	var (
		err    error
		result interface{}
	)
	if len(data) == 0 {
		return NewWithOption(nil, option), nil
	}
	if dataType == "" {
		dataType = checkDataType(data)
	}
	switch dataType {
	case "json", ".json", ".js":

	case "xml", ".xml":
		if data, err = gxml.ToJson(data); err != nil {
			return nil, err
		}

	case "yml", "yaml", ".yml", ".yaml":
		if data, err = gyaml.ToJson(data); err != nil {
			return nil, err
		}

	case "toml", ".toml":
		if data, err = gtoml.ToJson(data); err != nil {
			return nil, err
		}
	case "ini", ".ini":
		if data, err = gini.ToJson(data); err != nil {
			return nil, err
		}
	default:
		err = errors.New("unsupported type for loading")
	}
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if option.StrNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	switch result.(type) {
	case string, []byte:
		return nil, fmt.Errorf(`json decoding failed for content: %s`, string(data))
	}
	return NewWithOption(result, option), nil
}

// checkDataType automatically checks and returns the data type for <content>.
// Note that it uses regular expression for loose checking, you can use LoadXXX/LoadContentType
// functions to load the content for certain content type.
func checkDataType(content []byte) string {
	if json.Valid(content) {
		return "json"
	} else if gregex.IsMatch(`^<.+>[\S\s]+<.+>\s*$`, content) {
		return "xml"
	} else if !gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*"""[\s\S]+"""`, content) && !gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*'''[\s\S]+'''`, content) &&
		((gregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*".+"`, content) || gregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*\w+`, content)) ||
			(gregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*".+"`, content) || gregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*\w+`, content))) {
		return "yml"
	} else if !gregex.IsMatch(`^[\s\t\n\r]*;.+`, content) &&
		!gregex.IsMatch(`[\s\t\n\r]+;.+`, content) &&
		!gregex.IsMatch(`[\n\r]+[\s\t\w\-]+\.[\s\t\w\-]+\s*=\s*.+`, content) &&
		(gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, content) || gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, content)) {
		return "toml"
	} else if gregex.IsMatch(`\[[\w\.]+\]`, content) &&
		(gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, content) || gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, content)) {
		// Must contain "[xxx]" section.
		return "ini"
	} else {
		return ""
	}
}
