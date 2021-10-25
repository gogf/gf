// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gyaml provides accessing and converting for YAML content.
package gyaml

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"gopkg.in/yaml.v3"

	"github.com/gogf/gf/v2/util/gconv"
)

func Encode(value interface{}) (out []byte, err error) {
	if out, err = yaml.Marshal(value); err != nil {
		err = gerror.Wrap(err, `encode value to yaml failed`)
	}
	return
}

func Decode(value []byte) (interface{}, error) {
	var (
		result map[string]interface{}
		err    error
	)
	if err = yaml.Unmarshal(value, &result); err != nil {
		err = gerror.Wrap(err, `decode yaml failed`)
		return nil, err
	}
	return gconv.MapDeep(result), nil
}

func DecodeTo(value []byte, result interface{}) (err error) {
	err = yaml.Unmarshal(value, result)
	if err != nil {
		err = gerror.Wrap(err, `encode yaml to value failed`)
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
		if out, err = json.Marshal(result); err != nil {
			err = gerror.Wrap(err, `convert yaml to json failed`)
		}
		return out, err
	}
}
