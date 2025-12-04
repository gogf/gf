// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

// EnhancedStatusCode is http status for response.
type EnhancedStatusCode = int

// EnhancedStatusType is the structure for certain response status.
// Currently, it only supports `Response` and `Examples`.
// `Response` is the response structure
// `Examples` is the examples for the response, map[string]any, []any are supported.
type EnhancedStatusType struct {
	Response any
	Examples any
}

// IEnhanceResponseStatus is used to enhance the documentation of the response.
// Normal response structure could implement this interface to provide more information.
type IEnhanceResponseStatus interface {
	EnhanceResponseStatus() map[EnhancedStatusCode]EnhancedStatusType
}

// Response is specified by OpenAPI/Swagger 3.0 standard.
type Response struct {
	Description string      `json:"description"`
	Headers     Headers     `json:"headers,omitempty"`
	Content     Content     `json:"content,omitempty"`
	Links       Links       `json:"links,omitempty"`
	XExtensions XExtensions `json:"-"`
}

func (oai *OpenApiV3) tagMapToResponse(tagMap map[string]string, response *Response) error {
	var mergedTagMap = oai.fillMapWithShortTags(tagMap)
	if err := gconv.Struct(mergedTagMap, response); err != nil {
		return gerror.Wrap(err, `mapping struct tags to Response failed`)
	}
	oai.tagMapToXExtensions(mergedTagMap, response.XExtensions)
	return nil
}

func (r Response) MarshalJSON() ([]byte, error) {
	var (
		b   []byte
		m   map[string]json.RawMessage
		err error
	)
	type tempResponse Response // To prevent JSON marshal recursion error.
	if b, err = json.Marshal(tempResponse(r)); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	for k, v := range r.XExtensions {
		if b, err = json.Marshal(v); err != nil {
			return nil, err
		}
		m[k] = b
	}
	return json.Marshal(m)
}
