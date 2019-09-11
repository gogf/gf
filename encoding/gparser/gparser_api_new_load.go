// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

package gparser

import (
	"github.com/gogf/gf/encoding/gjson"
)

// New creates a Parser object with any variable type of <data>,
// but <data> should be a map or slice for data access reason,
// or it will make no sense.
// The <unsafe> param specifies whether using this Parser object
// in un-concurrent-safe context, which is false in default.
func New(value interface{}, safe ...bool) *Parser {
	return gjson.New(value, safe...)
}

// Load loads content from specified file <path>,
// and creates a Parser object from its content.
func Load(path string, safe ...bool) (*Parser, error) {
	return gjson.Load(path, safe...)
}

// LoadContent creates a Parser object from given content,
// it checks the data type of <content> automatically,
// supporting JSON, XML, YAML and TOML types of data.
func LoadContent(data interface{}, safe ...bool) (*Parser, error) {
	return gjson.LoadContent(data, safe...)
}

func LoadJson(data interface{}, safe ...bool) (*Parser, error) {
	return gjson.LoadJson(data, safe...)
}

func LoadXml(data interface{}, safe ...bool) (*Parser, error) {
	return gjson.LoadXml(data, safe...)
}

func LoadYaml(data interface{}, safe ...bool) (*Parser, error) {
	return gjson.LoadYaml(data, safe...)
}

func LoadToml(data interface{}, safe ...bool) (*Parser, error) {
	return gjson.LoadToml(data, safe...)
}

func LoadIni(data interface{}, safe ...bool) (*Parser, error) {
	return gjson.LoadIni(data, safe...)
}
