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
	return loadContentWithOptions(data, options)
}

// LoadJson creates a Json object from given JSON format content.
func LoadJson(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeJSON,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadXml creates a Json object from given XML format content.
func LoadXml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeXML,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadIni creates a Json object from given INI format content.
func LoadIni(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeIni,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadYaml creates a Json object from given YAML format content.
func LoadYaml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeYaml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadToml creates a Json object from given TOML format content.
func LoadToml(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeToml,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadProperties creates a Json object from given TOML format content.
func LoadProperties(data []byte, safe ...bool) (*Json, error) {
	var option = Options{
		Type: ContentTypeProperties,
	}
	if len(safe) > 0 && safe[0] {
		option.Safe = true
	}
	return loadContentWithOptions(data, option)
}

// LoadContent creates a Json object from given content, it checks the data type of `content`
// automatically, supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContent(data []byte, safe ...bool) (*Json, error) {
	return LoadContentType("", data, safe...)
}

// LoadContentType creates a Json object from given type and content,
// supporting data content type as follows:
// JSON, XML, INI, YAML and TOML.
func LoadContentType(dataType ContentType, data []byte, safe ...bool) (*Json, error) {
	if len(data) == 0 {
		return New(nil, safe...), nil
	}
	var options = Options{
		Type:      dataType,
		StrNumber: true,
	}
	if len(safe) > 0 && safe[0] {
		options.Safe = true
	}
	return loadContentWithOptions(data, options)
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
		ContentTypeJSON,
		ContentTypeJs,
		ContentTypeXML,
		ContentTypeYaml,
		ContentTypeYml,
		ContentTypeToml,
		ContentTypeIni,
		ContentTypeProperties:
		return true
	}
	return false
}

func trimBOM(data []byte) []byte {
	if len(data) < 3 {
		return data
	}
	if data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		data = data[3:]
	}
	return data
}

// loadContentWithOptions creates a Json object from given content.
// It supports data content type as follows:
// JSON, XML, INI, YAML and TOML.
func loadContentWithOptions(data []byte, options Options) (*Json, error) {
	var (
		err    error
		result any
	)
	data = trimBOM(data)
	if len(data) == 0 {
		return NewWithOptions(nil, options), nil
	}
	if options.Type == "" {
		options.Type, err = checkDataType(data)
		if err != nil {
			return nil, err
		}
	}
	options.Type = ContentType(gstr.TrimLeft(
		string(options.Type), "."),
	)
	switch options.Type {
	case ContentTypeJSON, ContentTypeJs:

	case ContentTypeXML:
		data, err = gxml.ToJson(data)

	case ContentTypeYaml, ContentTypeYml:
		data, err = gyaml.ToJson(data)

	case ContentTypeToml:
		data, err = gtoml.ToJson(data)

	case ContentTypeIni:
		data, err = gini.ToJson(data)

	case ContentTypeProperties:
		data, err = gproperties.ToJson(data)

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
// TODO it is not graceful here automatic judging the data type.
// TODO it might be removed in the future, which lets the user explicitly specify the data type not automatic checking.
func checkDataType(data []byte) (ContentType, error) {
	switch {
	case json.Valid(data):
		return ContentTypeJSON, nil

	case isXMLContent(data):
		return ContentTypeXML, nil

	case isYamlContent(data):
		return ContentTypeYaml, nil

	case isTomlContent(data):
		return ContentTypeToml, nil

	case isIniContent(data):
		// Must contain "[xxx]" section.
		return ContentTypeIni, nil

	case isPropertyContent(data):
		return ContentTypeProperties, nil

	default:
		return "", gerror.NewCode(
			gcode.CodeOperationFailed,
			`unable auto check the data format type`,
		)
	}
}

func isXMLContent(data []byte) bool {
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
