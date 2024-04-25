// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import (
	"bytes"

	"github.com/gogf/gf/v2/encoding/gini"
	"github.com/gogf/gf/v2/encoding/gproperties"
	"github.com/gogf/gf/v2/encoding/gtoml"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// LoadWithOptions creates a Json object from given JSON format content and options.
func LoadWithOptions(data []byte, options Options) (*Json, error) {
	return doLoadContentWithOptions(data, options)
}

// LoadJson creates a Json object from given JSON format content.
func LoadJson(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeJson,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(data, option)
}

// LoadXml creates a Json object from given XML format content.
func LoadXml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeXml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(data, option)
}

// LoadIni creates a Json object from given INI format content.
func LoadIni(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeIni,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(data, option)
}

// LoadYaml creates a Json object from given YAML format content.
func LoadYaml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeYaml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(data, option)
}

// LoadToml creates a Json object from given TOML format content.
func LoadToml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeToml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(data, option)
}

// LoadProperties creates a Json object from given TOML format content.
func LoadProperties(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeProperties,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return doLoadContentWithOptions(data, option)
}

// LoadContent creates a Json object from given content, it checks the data type of `content`
// automatically, supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContent(data []byte, safe ...bool) (*Json, error) {
	if len(data) == 0 {
		return New(nil, safe...), nil
	}
	return LoadContentType(checkDataType(data), data, safe...)
}

// LoadContentType creates a Json object from given type and content,
// supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContentType(dataType ContentType, data []byte, safe ...bool) (*Json, error) {
	if len(data) == 0 {
		return New(nil, safe...), nil
	}
	// ignore UTF8-BOM
	if data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	options := Options{
		Type:      dataType,
		StrNumber: true,
	}
	if len(safe) > 0 && safe[0] {
		options.Safe = true
	}
	return doLoadContentWithOptions(data, options)
}

// IsValidDataType checks and returns whether given `dataType` a valid data type for loading.
func IsValidDataType(dataType ContentType) bool {
	if dataType == "" {
		return false
	}
	if dataType[0] == '.' {
		dataType = dataType[1:]
	}
	switch dataType {
	case
		ContentTypeJson,
		ContentTypeJs,
		ContentTypeXml,
		ContentTypeYaml,
		ContentTypeYml,
		ContentTypeToml,
		ContentTypeIni,
		ContentTypeProperties:
		return true
	}
	return false
}

func loadContentWithOptions(data []byte, options Options) (*Json, error) {
	if len(data) == 0 {
		return NewWithOptions(nil, options), nil
	}
	if options.Type == "" {
		options.Type = checkDataType(data)
	}
	return loadContentTypeWithOptions(data, options)
}

func loadContentTypeWithOptions(data []byte, options Options) (*Json, error) {
	if len(data) == 0 {
		return NewWithOptions(nil, options), nil
	}
	// ignore UTF8-BOM
	if data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	return doLoadContentWithOptions(data, options)
}

// doLoadContent creates a Json object from given content.
// It supports data content type as follows:
// JSON, XML, INI, YAML and TOML.
func doLoadContentWithOptions(data []byte, options Options) (*Json, error) {
	var (
		err    error
		result interface{}
	)
	if len(data) == 0 {
		return NewWithOptions(nil, options), nil
	}
	if options.Type == "" {
		options.Type = checkDataType(data)
	}
	options.Type = ContentType(gstr.TrimLeft(
		string(options.Type), "."),
	)
	switch options.Type {
	case ContentTypeJson, ContentTypeJs:

	case ContentTypeXml:
		if data, err = gxml.ToJson(data); err != nil {
			return nil, err
		}

	case ContentTypeYaml, ContentTypeYml:
		if data, err = gyaml.ToJson(data); err != nil {
			return nil, err
		}

	case ContentTypeToml:
		if data, err = gtoml.ToJson(data); err != nil {
			return nil, err
		}

	case ContentTypeIni:
		if data, err = gini.ToJson(data); err != nil {
			return nil, err
		}
	case ContentTypeProperties:
		if data, err = gproperties.ToJson(data); err != nil {
			return nil, err
		}

	default:
		err = gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`unsupported type "%s" for loading`,
			options.Type,
		)
	}
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	if options.StrNumber {
		decoder.UseNumber()
	}
	if err = decoder.Decode(&result); err != nil {
		return nil, err
	}
	switch result.(type) {
	case string, []byte:
		return nil, gerror.Newf(`json decoding failed for content: %s`, data)
	}
	return NewWithOptions(result, options), nil
}

// checkDataType automatically checks and returns the data type for `content`.
// Note that it uses regular expression for loose checking, you can use LoadXXX/LoadContentType
// functions to load the content for certain content type.
func checkDataType(data []byte) ContentType {
	switch {
	case json.Valid(data):
		return ContentTypeJson
	case isXmlContent(data):
		return ContentTypeXml
	case isYamlContent(data):
		return ContentTypeYaml
	case isTomlContent(data):
		return ContentTypeToml
	case isIniContent(data):
		// Must contain "[xxx]" section.
		return ContentTypeIni
	case isPropertyContent(data):
		return ContentTypeProperties
	default:
		return ""
	}
}

func isXmlContent(data []byte) bool {
	return gregex.IsMatch(`^\s*<.+>[\S\s]+<.+>\s*$`, data)
}

func isYamlContent(data []byte) bool {
	return !gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*"""[\s\S]+"""`, data) &&
		!gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*'''[\s\S]+'''`, data) &&
		((gregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*".+"`, data) ||
			gregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*\w+`, data)) ||
			(gregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*".+"`, data) ||
				gregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*\w+`, data)))
}

func isTomlContent(data []byte) bool {
	return !gregex.IsMatch(`^[\s\t\n\r]*;.+`, data) &&
		!gregex.IsMatch(`[\s\t\n\r]+;.+`, data) &&
		!gregex.IsMatch(`[\n\r]+[\s\t\w\-]+\.[\s\t\w\-]+\s*=\s*.+`, data) &&
		(gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, data) ||
			gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data))
}

func isIniContent(data []byte) bool {
	return gregex.IsMatch(`\[[\w\.]+\]`, data) &&
		(gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, data) ||
			gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data))
}

func isPropertyContent(data []byte) bool {
	return gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data)
}
