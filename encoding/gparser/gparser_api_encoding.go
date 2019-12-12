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
	return New(value).ToXml(rootTag...)
}

func VarToXmlString(value interface{}, rootTag ...string) (string, error) {
	return New(value).ToXmlString(rootTag...)
}

func VarToXmlIndent(value interface{}, rootTag ...string) ([]byte, error) {
	return New(value).ToXmlIndent(rootTag...)
}

func VarToXmlIndentString(value interface{}, rootTag ...string) (string, error) {
	return New(value).ToXmlIndentString(rootTag...)
}

func MustToXml(value interface{}, rootTag ...string) []byte {
	return New(value).MustToXml(rootTag...)
}

func MustToXmlString(value interface{}, rootTag ...string) string {
	return New(value).MustToXmlString(rootTag...)
}

func MustToXmlIndent(value interface{}, rootTag ...string) []byte {
	return New(value).MustToXmlIndent(rootTag...)
}

func MustToXmlIndentString(value interface{}, rootTag ...string) string {
	return New(value).MustToXmlIndentString(rootTag...)
}

// ========================================================================
// YAML
// ========================================================================

func VarToYaml(value interface{}) ([]byte, error) {
	return New(value).ToYaml()
}

func VarToYamlString(value interface{}) (string, error) {
	return New(value).ToYamlString()
}

func MustToYaml(value interface{}) []byte {
	return New(value).MustToYaml()
}

func MustToYamlString(value interface{}) string {
	return New(value).MustToYamlString()
}

// ========================================================================
// TOML
// ========================================================================

func VarToToml(value interface{}) ([]byte, error) {
	return New(value).ToToml()
}

func VarToTomlString(value interface{}) (string, error) {
	return New(value).ToTomlString()
}

func MustToToml(value interface{}) []byte {
	return New(value).MustToToml()
}

func MustToTomlString(value interface{}) string {
	return New(value).MustToTomlString()
}

// ========================================================================
// INI
// ========================================================================

func VarToIni(value interface{}) ([]byte, error) {
	return New(value).ToIni()
}

func VarToIniString(value interface{}) (string, error) {
	return New(value).ToIniString()
}

func MustToIni(value interface{}) []byte {
	return New(value).MustToIni()
}

func MustToIniString(value interface{}) string {
	return New(value).MustToIniString()
}
