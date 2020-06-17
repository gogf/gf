// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"github.com/gogf/gf/encoding/gini"
	"github.com/gogf/gf/encoding/gtoml"
	"github.com/gogf/gf/encoding/gxml"
	"github.com/gogf/gf/encoding/gyaml"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/util/gconv"
)

// ========================================================================
// JSON
// ========================================================================

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

func (j *Json) MustToJson() []byte {
	result, err := j.ToJson()
	if err != nil {
		panic(err)
	}
	return result
}

func (j *Json) MustToJsonString() string {
	return gconv.UnsafeBytesToStr(j.MustToJson())
}

func (j *Json) MustToJsonIndent() []byte {
	result, err := j.ToJsonIndent()
	if err != nil {
		panic(err)
	}
	return result
}

func (j *Json) MustToJsonIndentString() string {
	return gconv.UnsafeBytesToStr(j.MustToJsonIndent())
}

// ========================================================================
// XML
// ========================================================================

func (j *Json) ToXml(rootTag ...string) ([]byte, error) {
	return gxml.Encode(j.ToMap(), rootTag...)
}

func (j *Json) ToXmlString(rootTag ...string) (string, error) {
	b, e := j.ToXml(rootTag...)
	return string(b), e
}

func (j *Json) ToXmlIndent(rootTag ...string) ([]byte, error) {
	return gxml.EncodeWithIndent(j.ToMap(), rootTag...)
}

func (j *Json) ToXmlIndentString(rootTag ...string) (string, error) {
	b, e := j.ToXmlIndent(rootTag...)
	return string(b), e
}

func (j *Json) MustToXml(rootTag ...string) []byte {
	result, err := j.ToXml(rootTag...)
	if err != nil {
		panic(err)
	}
	return result
}

func (j *Json) MustToXmlString(rootTag ...string) string {
	return gconv.UnsafeBytesToStr(j.MustToXml(rootTag...))
}

func (j *Json) MustToXmlIndent(rootTag ...string) []byte {
	result, err := j.ToXmlIndent(rootTag...)
	if err != nil {
		panic(err)
	}
	return result
}

func (j *Json) MustToXmlIndentString(rootTag ...string) string {
	return gconv.UnsafeBytesToStr(j.MustToXmlIndent(rootTag...))
}

// ========================================================================
// YAML
// ========================================================================

func (j *Json) ToYaml() ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gyaml.Encode(*(j.p))
}

func (j *Json) ToYamlString() (string, error) {
	b, e := j.ToYaml()
	return string(b), e
}

func (j *Json) MustToYaml() []byte {
	result, err := j.ToYaml()
	if err != nil {
		panic(err)
	}
	return result
}

func (j *Json) MustToYamlString() string {
	return gconv.UnsafeBytesToStr(j.MustToYaml())
}

// ========================================================================
// TOML
// ========================================================================

func (j *Json) ToToml() ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gtoml.Encode(*(j.p))
}

func (j *Json) ToTomlString() (string, error) {
	b, e := j.ToToml()
	return string(b), e
}

func (j *Json) MustToToml() []byte {
	result, err := j.ToToml()
	if err != nil {
		panic(err)
	}
	return result
}

func (j *Json) MustToTomlString() string {
	return gconv.UnsafeBytesToStr(j.MustToToml())
}

// ========================================================================
// INI
// ========================================================================

func (j *Json) ToIni() ([]byte, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return gini.Encode((*(j.p)).(map[string]interface{}))
}

func (j *Json) ToIniString() (string, error) {
	b, e := j.ToToml()
	return string(b), e
}

func (j *Json) MustToIni() []byte {
	result, err := j.ToIni()
	if err != nil {
		panic(err)
	}
	return result
}

func (j *Json) MustToIniString() string {
	return gconv.UnsafeBytesToStr(j.MustToIni())
}
