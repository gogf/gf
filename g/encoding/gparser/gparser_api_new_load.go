// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

package gparser

import (
	"github.com/gogf/gf/g/encoding/gjson"
)

// New creates a Parser object with any variable type of <data>,
// but <data> should be a map or slice for data access reason,
// or it will make no sense.
// The <unsafe> param specifies whether using this Parser object
// in un-concurrent-safe context, which is false in default.
func New(value interface{}, unsafe...bool) *Parser {
    return &Parser{gjson.New(value, unsafe...)}
}

// NewUnsafe creates a un-concurrent-safe Parser object.
func NewUnsafe(value...interface{}) *Parser {
    if len(value) > 0 {
        return &Parser{gjson.New(value[0], false)}
    }
    return &Parser{gjson.New(nil, false)}
}

// Load loads content from specified file <path>,
// and creates a Parser object from its content.
func Load(path string, unsafe...bool) (*Parser, error) {
    if j, e := gjson.Load(path, unsafe...); e == nil {
        return &Parser{j}, nil
    } else {
        return nil, e
    }
}

// LoadContent creates a Parser object from given content,
// it checks the data type of <content> automatically,
// supporting JSON, XML, YAML and TOML types of data.
func LoadContent(data []byte, unsafe...bool) (*Parser, error) {
    if j, e := gjson.LoadContent(data, unsafe...); e == nil {
        return &Parser{j}, nil
    } else {
        return nil, e
    }
}
