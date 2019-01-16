// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gyaml provides accessing and converting for YAML content.
package gyaml

import "gitee.com/johng/gf/third/github.com/ghodss/yaml"

func Encode(v interface{}) ([]byte, error) {
    return yaml.Marshal(v)
}

func Decode(v []byte) (interface{}, error) {
    var result interface{}
    if err := yaml.Unmarshal(v, &result); err != nil {
        return nil, err
    }
    return result, nil
}

func DecodeTo(v []byte, result interface{}) error {
    return yaml.Unmarshal(v, &result)
}

func ToJson(v []byte) ([]byte, error) {
    return yaml.YAMLToJSON(v)
}