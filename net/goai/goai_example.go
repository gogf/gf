// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"fmt"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
)

// Example is specified by OpenAPI/Swagger 3.0 standard.
type Example struct {
	Summary       string `json:"summary,omitempty"`
	Description   string `json:"description,omitempty"`
	Value         any    `json:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty"`
}

type Examples map[string]*ExampleRef

type ExampleRef struct {
	Ref   string
	Value *Example
}

func (e *Examples) applyExamplesFile(path string) error {
	if empty.IsNil(e) {
		return nil
	}
	var json string
	if resource := gres.Get(path); resource != nil {
		json = string(resource.Content())
	} else {
		absolutePath := gfile.RealPath(path)
		if absolutePath != "" {
			json = gfile.GetContents(absolutePath)
		}
	}
	if json == "" {
		return nil
	}
	var data any
	err := gjson.Unmarshal([]byte(json), &data)
	if err != nil {
		return err
	}
	err = e.applyExamplesData(data)
	if err != nil {
		return err
	}
	return nil
}

func (e *Examples) applyExamplesData(data any) error {
	if empty.IsNil(e) || empty.IsNil(data) {
		return nil
	}

	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			(*e)[key] = &ExampleRef{
				Value: &Example{
					Value: value,
				},
			}
		}
	case []any:
		for i, value := range v {
			(*e)[fmt.Sprintf("example %d", i+1)] = &ExampleRef{
				Value: &Example{
					Value: value,
				},
			}
		}
	default:
	}
	return nil
}

func (r ExampleRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}
