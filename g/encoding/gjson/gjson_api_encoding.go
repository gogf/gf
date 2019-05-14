// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"encoding/json"
	"github.com/gogf/gf/g/encoding/gtoml"
	"github.com/gogf/gf/g/encoding/gxml"
	"github.com/gogf/gf/g/encoding/gyaml"
)

func (j *Json) ToXml(rootTag...string) ([]byte, error) {
    return gxml.Encode(j.ToMap(), rootTag...)
}

func (j *Json) ToXmlString(rootTag...string) (string, error) {
    b, e := j.ToXml(rootTag...)
    return string(b), e
}

func (j *Json) ToXmlIndent(rootTag...string) ([]byte, error) {
    return gxml.EncodeWithIndent(j.ToMap(), rootTag...)
}

func (j *Json) ToXmlIndentString(rootTag...string) (string, error) {
    b, e := j.ToXmlIndent(rootTag...)
    return string(b), e
}

func (j *Json) ToJson() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return Encode(*(j.p))
}

func (j *Json) ToJsonString() (string, error) {
    b, e := j.ToJson()
    return string(b), e
}

func (j *Json) ToJsonIndent() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return json.MarshalIndent(*(j.p), "", "\t")
}

func (j *Json) ToJsonIndentString() (string, error) {
    b, e := j.ToJsonIndent()
    return string(b), e
}

func (j *Json) ToYaml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gyaml.Encode(*(j.p))
}

func (j *Json) ToYamlString() (string, error) {
    b, e := j.ToYaml()
    return string(b), e
}

func (j *Json) ToToml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gtoml.Encode(*(j.p))
}

func (j *Json) ToTomlString() (string, error) {
    b, e := j.ToToml()
    return string(b), e
}
