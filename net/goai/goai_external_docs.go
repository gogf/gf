// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// ExternalDocs is specified by OpenAPI/Swagger standard version 3.0.
type ExternalDocs struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

func (ed *ExternalDocs) UnmarshalValue(value interface{}) error {
	var valueBytes = gconv.Bytes(value)
	if json.Valid(valueBytes) {
		return json.UnmarshalUseNumber(valueBytes, ed)
	}
	var (
		valueString = string(valueBytes)
		valueArray  = gstr.Split(valueString, "|")
	)
	ed.URL = valueArray[0]
	if len(valueArray) > 1 {
		ed.Description = valueArray[1]
	}
	return nil
}
