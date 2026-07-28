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
	var (
		checkType   ContentType
		decodedData any
	)
	if options.Type != "" {
		checkType = gstr.TrimLeft(options.Type, ".")
	} else {
		checkType, err = checkDataType(data)
		if err != nil {
			return nil, err
		}
	}
	switch checkType {
	case ContentTypeJSON, ContentTypeJs:
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

	case ContentTypeXML:
		decodedData, err = gxml.Decode(data)
		if err != nil {
			return nil, err
		}
		return NewWithOptions(decodedData, options), nil

	case ContentTypeYaml, ContentTypeYml:
		decodedData, err = gyaml.Decode(data)
		if err != nil {
			return nil, err
		}
		return NewWithOptions(decodedData, options), nil

	case ContentTypeToml:
		decodedData, err = gtoml.Decode(data)
		if err != nil {
			return nil, err
		}
		return NewWithOptions(decodedData, options), nil

	case ContentTypeIni:
		decodedData, err = gini.Decode(data)
		if err != nil {
			return nil, err
		}
		return NewWithOptions(decodedData, options), nil

	case ContentTypeProperties:
		decodedData, err = gproperties.Decode(data)
		if err != nil {
			return nil, err
		}
		return NewWithOptions(decodedData, options), nil

	default:
	}
	// ignore some duplicated types, like js and yml,
	// which are not necessary shown in error message.
	allSupportedTypes := []string{
		ContentTypeJSON,
		ContentTypeXML,
		ContentTypeYaml,
		ContentTypeToml,
		ContentTypeIni,
		ContentTypeProperties,
	}
	return nil, gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`unsupported type "%s" for loading, all supported types: %s`,
		options.Type, gstr.Join(allSupportedTypes, ", "),
	)
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

// isXMLContent checks whether given content is XML format.
// XML format is easy to be identified using regular expression.
func isXMLContent(data []byte) bool {
	return gregex.IsMatch(`^\s*<.+>[\S\s]+<.+>\s*$`, data)
}

// isYamlContent checks whether given content is YAML format.
func isYamlContent(data []byte) bool {
	// x = y
	// "x.x" = "y"
	tomlFormat1 := gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*"""[\s\S]+"""`, data)
	if tomlFormat1 {
		return false
	}
	// "x.x" = '''
	// y
	// '''
	tomlFormat2 := gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*'''[\s\S]+'''`, data)
	if tomlFormat2 {
		return false
	}

	// content starts with:
	// x : "y"
	yamlFormat1 := gregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s+".+"`, data)

	// content starts with:
	// x : y
	yamlFormat2 := gregex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s+\w+`, data)

	// line starts with:
	// x : "y"
	yamlFormat3 := gregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s+".+"`, data)

	// line starts with:
	// x : y
	yamlFormat4 := gregex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s+\w+`, data)

	// content starts with:
	// "x" : "y"
	yamlFormat5 := gregex.IsMatch(`^[\n\r]*".+":\s+".+"`, data)

	// line starts with:
	// "x" : y
	yamlFormat6 := gregex.IsMatch(`[\n\r]+".+":\s+\w+`, data)

	return yamlFormat1 || yamlFormat2 || yamlFormat3 || yamlFormat4 || yamlFormat5 || yamlFormat6
}

// isTomlContent checks whether given content is TOML format.
func isTomlContent(data []byte) bool {
	// content starts with:
	// ; comment line
	contentStartsWithSemicolonComment := gregex.IsMatch(`^[\s\t\n\r]*;.+`, data)
	if contentStartsWithSemicolonComment {
		return false
	}
	// line starts with:
	// ; comment line
	lineStartsWithSemicolonComment := gregex.IsMatch(`[\s\t\n\r]+;.+`, data)
	if lineStartsWithSemicolonComment {
		return false
	}

	// line starts with, this should not be toml format:
	// key.with.dot = value
	keyWithDot := gregex.IsMatch(`[\n\r]+[\s\t\w\-]+\.[\s\t\w\-]+\s*=\s*.+`, data)
	if keyWithDot {
		return false
	}

	// line starts with:
	// key = value
	// key = "value"
	// "key" = "value"
	// "key" = value
	tomlFormat1 := gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, data)
	tomlFormat2 := gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data)
	return tomlFormat1 || tomlFormat2
}

// isIniContent checks whether given content is INI format.
func isIniContent(data []byte) bool {
	// no section like: [section], but ini format must have sections.
	hasBrackets := gregex.IsMatch(`\[[\w\.]+\]`, data)
	if !hasBrackets {
		return false
	}
	iniFormat1 := gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, data)
	iniFormat2 := gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data)
	return iniFormat1 || iniFormat2
}

// isPropertyContent checks whether given content is Properties format.
func isPropertyContent(data []byte) bool {
	// line starts with:
	// key = value
	// "key" = value
	propertyFormat := gregex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, data)
	return propertyFormat
}
