// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gjson provides convenient API for JSON/XML/YAML/TOML data handling.
package gjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/gogf/gf/encoding/gini"
	"github.com/gogf/gf/encoding/gtoml"
	"github.com/gogf/gf/encoding/gxml"
	"github.com/gogf/gf/encoding/gyaml"
	"github.com/gogf/gf/internal/rwmutex"
	"github.com/gogf/gf/os/gfcache"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// New creates a Json object with any variable type of <data>, but <data> should be a map or slice
// for data access reason, or it will make no sense.
//
// The parameter <safe> specifies whether using this Json object in concurrent-safe context, which is
// false in default.
func New(data interface{}, safe ...bool) *Json {
	return NewWithTag(data, "json", safe...)
}

// NewWithTag creates a Json object with any variable type of <data>, but <data> should be a map or slice
// for data access reason, or it will make no sense.
//
// The parameter <tags> specifies priority tags for struct conversion to map, multiple tags joined with char ','.
//
// The parameter <safe> specifies whether using this Json object in concurrent-safe context, which is
// false in default.
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
			// Note that it uses MapDeep function implementing the converting.
			i = gconv.MapDeep(data, tags)
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

// Valid checks whether <data> is a valid JSON data type.
func Valid(data interface{}) bool {
	return json.Valid(gconv.Bytes(data))
}

// Encode encodes <value> to JSON data type of bytes.
func Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Decode decodes <data>(string/[]byte) to golang variable.
func Decode(data interface{}) (interface{}, error) {
	var value interface{}
	if err := DecodeTo(gconv.Bytes(data), &value); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

// Decode decodes <data>(string/[]byte) to specified golang variable <v>.
// The <v> should be a pointer type.
func DecodeTo(data interface{}, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(gconv.Bytes(data)))
	// Do not use number, it converts float64 to json.Number type,
	// which actually a string type. It causes converting issue for other data formats,
	// for example: yaml.
	//decoder.UseNumber()
	return decoder.Decode(v)
}

// DecodeToJson codes <data>(string/[]byte) to a Json object.
func DecodeToJson(data interface{}, safe ...bool) (*Json, error) {
	if v, err := Decode(gconv.Bytes(data)); err != nil {
		return nil, err
	} else {
		return New(v, safe...), nil
	}
}

// Load loads content from specified file <path>,
// and creates a Json object from its content.
func Load(path string, safe ...bool) (*Json, error) {
	if p, err := gfile.Search(path); err != nil {
		return nil, err
	} else {
		path = p
	}
	return doLoadContent(gfile.Ext(path), gfcache.GetBinContents(path), safe...)
}

func LoadJson(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("json", gconv.Bytes(data), safe...)
}

func LoadXml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("xml", gconv.Bytes(data), safe...)
}

func LoadYaml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("yaml", gconv.Bytes(data), safe...)
}

func LoadToml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("toml", gconv.Bytes(data), safe...)
}

func LoadIni(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("ini", gconv.Bytes(data), safe...)
}

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

func LoadContent(data interface{}, safe ...bool) (*Json, error) {
	content := gconv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}
	return doLoadContent(checkDataType(content), content, safe...)

}

// checkDataType automatically checks and returns the data type for <content>.
func checkDataType(content []byte) string {
	if json.Valid(content) {
		return "json"
	} else if gregex.IsMatch(`^<.+>[\S\s]+<.+>$`, content) {
		return "xml"
	} else if gregex.IsMatch(`^[\s\t]*[\w\-]+\s*:\s*.+`, content) || gregex.IsMatch(`\n[\s\t]*[\w\-]+\s*:\s*.+`, content) {
		return "yml"
	} else if (gregex.IsMatch(`^[\s\t\[*\]].?*[\w\-]+\s*=\s*.+`, content) || gregex.IsMatch(`\n[\s\t\[*\]]*[\w\-]+\s*=\s*.+`, content)) && gregex.IsMatch(`\n[\s\t]*[\w\-]+\s*=*\"*.+\"`, content) == false && gregex.IsMatch(`^[\s\t]*[\w\-]+\s*=*\"*.+\"`, content) == false {
		return "ini"
	} else if gregex.IsMatch(`^[\s\t]*[\w\-]+\s*=\s*.+`, content) || gregex.IsMatch(`\n[\s\t]*[\w\-]+\s*=\s*.+`, content) {
		return "toml"
	} else {
		return ""
	}
}
