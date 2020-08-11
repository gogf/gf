// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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
	j := (*Json)(nil)
	switch data.(type) {
	case string, []byte:
		if r, err := LoadContent(gconv.Bytes(data)); err == nil {
			j = r
		} else {
			j = &Json{
				p:  &data,
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		}
	default:
		rv := reflect.ValueOf(data)
		kind := rv.Kind()
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
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		case reflect.Map, reflect.Struct:
			i := interface{}(nil)
			// Note that it uses Map function implementing the converting.
			// Note that it here should not use MapDeep function if you really know what it means.
			i = gconv.Map(data, tags)
			j = &Json{
				p:  &i,
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		default:
			j = &Json{
				p:  &data,
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		}
	}
	j.mu = rwmutex.New(safe...)
	return j
}

// Load loads content from specified file <path>, and creates a Json object from its content.
func Load(path string, safe ...bool) (*Json, error) {
	if p, err := gfile.Search(path); err != nil {
		return nil, err
	} else {
		path = p
	}
	return doLoadContent(gfile.Ext(path), gfile.GetBytesWithCache(path), safe...)
}

// LoadJson creates a Json object from given JSON format content.
func LoadJson(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("json", gconv.Bytes(data), safe...)
}

// LoadXml creates a Json object from given XML format content.
func LoadXml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("xml", gconv.Bytes(data), safe...)
}

// LoadIni creates a Json object from given INI format content.
func LoadIni(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("ini", gconv.Bytes(data), safe...)
}

// LoadYaml creates a Json object from given YAML format content.
func LoadYaml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("yaml", gconv.Bytes(data), safe...)
}

// LoadToml creates a Json object from given TOML format content.
func LoadToml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("toml", gconv.Bytes(data), safe...)
}

// doLoadContent creates a Json object from given content.
// It supports data content type as follows:
// JSON, XML, INI, YAML and TOML.
func doLoadContent(dataType string, data []byte, safe ...bool) (*Json, error) {
	var err error
	var result interface{}
	if len(data) == 0 {
		return New(nil, safe...), nil
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
	// Do not use number, it converts float64 to json.Number type,
	// which actually a string type. It causes converting issue for other data formats,
	// for example: yaml.
	//decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	switch result.(type) {
	case string, []byte:
		return nil, fmt.Errorf(`json decoding failed for content: %s`, string(data))
	}
	return New(result, safe...), nil
}

// LoadContent creates a Json object from given content, it checks the data type of <content>
// automatically, supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContent(data interface{}, safe ...bool) (*Json, error) {
	content := gconv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}

	//ignore UTF8-BOM
	if content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}

	return doLoadContent(checkDataType(content), content, safe...)

}

// checkDataType automatically checks and returns the data type for <content>.
// Note that it uses regular expression for loose checking, you can use LoadXXX
// functions to load the content for certain content type.
func checkDataType(content []byte) string {
	if json.Valid(content) {
		return "json"
	} else if gregex.IsMatch(`^<.+>[\S\s]+<.+>$`, content) {
		return "xml"
	} else if gregex.IsMatch(`[\s\t\n]*[\w\-]+\s*:\s*.+`, content) {
		return "yml"
	} else if gregex.IsMatch(`\[[\w]+\]`, content) &&
		gregex.IsMatch(`[\s\t\n\[\]]*[\w\-]+\s*=\s*.+`, content) &&
		!gregex.IsMatch(`[\s\t\n]*[\w\-]+\s*=*\"*.+\"`, content) {
		// Must contain "[xxx]" section.
		return "ini"
	} else if gregex.IsMatch(`[\s\t\n]*[\w\-\."]+\s*=\s*.+`, content) {
		return "toml"
	} else {
		return ""
	}
}
