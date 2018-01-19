// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// YAML
package gyaml

import "github.com/ghodss/yaml"

func Encode(v interface{}) ([]byte, error) {
    return yaml.Marshal(v)
}

func Decode(v []byte) error {
    var result interface{}
    return yaml.Unmarshal(v, &result)
}

func DecodeTo(v []byte, result interface{}) error {
    return yaml.Unmarshal(v, &result)
}

func ToJson(v []byte) ([]byte, error) {
    return yaml.YAMLToJSON(v)
}