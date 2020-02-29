// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

package gparser

// ========================================================================
// JSON
// ========================================================================

func VarToJson(value interface{}) ([]byte, error) {
	return New(value).ToJson()
}

func VarToJsonString(value interface{}) (string, error) {
	return New(value).ToJsonString()
}

func VarToJsonIndent(value interface{}) ([]byte, error) {
	return New(value).ToJsonIndent()
}

func VarToJsonIndentString(value interface{}) (string, error) {
	return New(value).ToJsonIndentString()
}

func MustToJson(value interface{}) []byte {
	return New(value).MustToJson()
}

func MustToJsonString(value interface{}) string {
	return New(value).MustToJsonString()
}

func MustToJsonIndent(value interface{}) []byte {
	return New(value).MustToJsonIndent()
}

func MustToJsonIndentString(value interface{}) string {
	return New(value).MustToJsonIndentString()
}

// ========================================================================
// XML
// ========================================================================

func VarToXml(value interface{}, rootTag ...string) ([]byte, error) {
	return NewWithTag(value, "xml").ToXml(rootTag...)
}

func VarToXmlString(value interface{}, rootTag ...string) (string, error) {
	return NewWithTag(value, "xml").ToXmlString(rootTag...)
}

func VarToXmlIndent(value interface{}, rootTag ...string) ([]byte, error) {
	return NewWithTag(value, "xml").ToXmlIndent(rootTag...)
}

func VarToXmlIndentString(value interface{}, rootTag ...string) (string, error) {
	return NewWithTag(value, "xml").ToXmlIndentString(rootTag...)
}

func MustToXml(value interface{}, rootTag ...string) []byte {
	return NewWithTag(value, "xml").MustToXml(rootTag...)
}

func MustToXmlString(value interface{}, rootTag ...string) string {
	return NewWithTag(value, "xml").MustToXmlString(rootTag...)
}

func MustToXmlIndent(value interface{}, rootTag ...string) []byte {
	return NewWithTag(value, "xml").MustToXmlIndent(rootTag...)
}

func MustToXmlIndentString(value interface{}, rootTag ...string) string {
	return NewWithTag(value, "xml").MustToXmlIndentString(rootTag...)
}

// ========================================================================
// YAML
// ========================================================================

func VarToYaml(value interface{}) ([]byte, error) {
	return NewWithTag(value, "yaml").ToYaml()
}

func VarToYamlString(value interface{}) (string, error) {
	return NewWithTag(value, "yaml").ToYamlString()
}

func MustToYaml(value interface{}) []byte {
	return NewWithTag(value, "yaml").MustToYaml()
}

func MustToYamlString(value interface{}) string {
	return NewWithTag(value, "yaml").MustToYamlString()
}

// ========================================================================
// TOML
// ========================================================================

func VarToToml(value interface{}) ([]byte, error) {
	return NewWithTag(value, "toml").ToToml()
}

func VarToTomlString(value interface{}) (string, error) {
	return NewWithTag(value, "toml").ToTomlString()
}

func MustToToml(value interface{}) []byte {
	return NewWithTag(value, "toml").MustToToml()
}

func MustToTomlString(value interface{}) string {
	return NewWithTag(value, "toml").MustToTomlString()
}

// ========================================================================
// INI
// ========================================================================

func VarToIni(value interface{}) ([]byte, error) {
	return NewWithTag(value, "ini").ToIni()
}

func VarToIniString(value interface{}) (string, error) {
	return NewWithTag(value, "ini").ToIniString()
}

func MustToIni(value interface{}) []byte {
	return NewWithTag(value, "ini").MustToIni()
}

func MustToIniString(value interface{}) string {
	return NewWithTag(value, "ini").MustToIniString()
}
