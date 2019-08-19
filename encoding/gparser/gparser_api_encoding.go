// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

package gparser

func VarToXml(value interface{}, rootTag ...string) ([]byte, error) {
	return New(value).ToXml(rootTag...)
}

func VarToXmlIndent(value interface{}, rootTag ...string) ([]byte, error) {
	return New(value).ToXmlIndent(rootTag...)
}

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

func VarToYaml(value interface{}) ([]byte, error) {
	return New(value).ToYaml()
}

func VarToToml(value interface{}) ([]byte, error) {
	return New(value).ToToml()
}

func VarToIni(value interface{}) ([]byte, error) {
	return New(value).ToIni()
}
