// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gyaml provides accessing and converting for YAML content.
package gyaml

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

func Encode(value interface{}) (out []byte, err error) {
	if out, err = yaml.Marshal(value); err != nil {
		err = gerror.Wrap(err, `yaml.Marshal failed`)
	}
	return
}

func EncodeIndent(value interface{}, indent string) (out []byte, err error) {
	out, err = Encode(value)
	if err != nil {
		return
	}
	if indent != "" {
		var (
			buffer = bytes.NewBuffer(nil)
			array  = strings.Split(strings.TrimSpace(string(out)), "\n")
		)
		for _, v := range array {
			buffer.WriteString(indent)
			buffer.WriteString(v)
			buffer.WriteString("\n")
		}
		out = buffer.Bytes()
	}
	return
}

func Decode(value []byte) (interface{}, error) {
	var (
		result map[string]interface{}
		err    error
	)
	if err = yaml.Unmarshal(value, &result); err != nil {
		err = gerror.Wrap(err, `yaml.Unmarshal failed`)
		return nil, err
	}
	return gconv.MapDeep(result), nil
}

func DecodeTo(value []byte, result interface{}) (err error) {
	err = yaml.Unmarshal(value, result)
	if err != nil {
		err = gerror.Wrap(err, `yaml.Unmarshal failed`)
	}
	return
}

func ToJson(value []byte) (out []byte, err error) {
	var (
		result interface{}
	)
	if result, err = Decode(value); err != nil {
		return nil, err
	} else {
		return json.Marshal(result)
	}
}
