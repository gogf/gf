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
	"github.com/gogf/gf/g/encoding/gtoml"
	"github.com/gogf/gf/g/encoding/gxml"
	"github.com/gogf/gf/g/encoding/gyaml"
	"github.com/gogf/gf/g/internal/rwmutex"
	"github.com/gogf/gf/g/os/gfcache"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/util/gconv"
	"reflect"
)

// New creates a Json object with any variable type of <data>,
// but <data> should be a map or slice for data access reason,
// or it will make no sense.
// The <unsafe> param specifies whether using this Json object
// in un-concurrent-safe context, which is false in default.
func New(data interface{}, unsafe...bool) *Json {
    j := (*Json)(nil)
    switch data.(type) {
        case string, []byte:
            if r, err := LoadContent(gconv.Bytes(data)); err == nil {
	            j = r
            } else {
	            j = &Json {
		            p  : &data,
		            c  : byte(gDEFAULT_SPLIT_CHAR),
		            vc : false ,
	            }
            }
        default:
	        rv   := reflect.ValueOf(data)
	        kind := rv.Kind()
	        switch kind {
		        case reflect.Slice: fallthrough
		        case reflect.Array:
			        i := interface{}(nil)
			        i  = gconv.Interfaces(data)
			        j  = &Json {
				        p  : &i,
				        c  : byte(gDEFAULT_SPLIT_CHAR),
				        vc : false ,
			        }
		        case reflect.Map: fallthrough
		        case reflect.Struct:
			        i := interface{}(nil)
			        i  = gconv.Map(data)
			        j  = &Json {
				        p  : &i,
				        c  : byte(gDEFAULT_SPLIT_CHAR),
				        vc : false ,
			        }
		        default:
			        j  = &Json {
				        p  : &data,
				        c  : byte(gDEFAULT_SPLIT_CHAR),
				        vc : false ,
			        }
	        }
    }
    j.mu = rwmutex.New(unsafe...)
    return j
}

// NewUnsafe creates a un-concurrent-safe Json object.
func NewUnsafe(data...interface{}) *Json {
    if len(data) > 0 {
        return New(data[0], true)
    }
    return New(nil, true)
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
    decoder.UseNumber()
    return decoder.Decode(v)
}

// DecodeToJson codes <data>(string/[]byte) to a Json object.
func DecodeToJson(data interface{}, unsafe...bool) (*Json, error) {
    if v, err := Decode(gconv.Bytes(data)); err != nil {
        return nil, err
    } else {
        return New(v, unsafe...), nil
    }
}

// Load loads content from specified file <path>,
// and creates a Json object from its content.
func Load(path string, unsafe...bool) (*Json, error) {
    return LoadContent(gfcache.GetBinContents(path), unsafe...)
}

// LoadContent creates a Json object from given content,
// it checks the data type of <content> automatically,
// supporting JSON, XML, YAML and TOML types of data.
func LoadContent(data interface{}, unsafe...bool) (*Json, error) {
    var err    error
    var result interface{}
    b := gconv.Bytes(data)
    t := "json"
    // auto check data type
    if json.Valid(b) {
        t = "json"
    } else if gregex.IsMatch(`^<.+>.*</.+>$`, b) {
        t = "xml"
    } else if gregex.IsMatch(`^[\s\t]*\w+\s*:\s*.+`, b) || gregex.IsMatch(`\n[\s\t]*\w+\s*:\s*.+`, b) {
        t = "yml"
    } else if gregex.IsMatch(`^[\s\t]*\w+\s*=\s*.+`, b) || gregex.IsMatch(`\n[\s\t]*\w+\s*=\s*.+`, b) {
        t = "toml"
    } else {
        return nil, errors.New("unsupported data type")
    }
    // convert to json type data
    switch t {
        case "json", ".json":
            // ok
        case "xml", ".xml":
            // TODO UseNumber
            b, err = gxml.ToJson(b)

        case "yml", "yaml", ".yml", ".yaml":
            // TODO UseNumber
            b, err = gyaml.ToJson(b)

        case "toml", ".toml":
            // TODO UseNumber
            b, err = gtoml.ToJson(b)

        default:
            err = errors.New("nonsupport type " + t)
    }
    if err != nil {
        return nil, err
    }
    if result == nil {
        decoder := json.NewDecoder(bytes.NewReader(b))
        decoder.UseNumber()
        if err := decoder.Decode(&result); err != nil {
            return nil, err
        }
        switch result.(type) {
            case string, []byte:
                return nil, fmt.Errorf(`json decoding failed for content: %s`, string(b))
        }
    }
    return New(result, unsafe...), nil
}
