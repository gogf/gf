// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gxml provides accessing and converting for XML content.
package gxml

import (
	"strings"

	"github.com/clbanning/mxj"
	"github.com/gogf/gf/encoding/gcharset"
	"github.com/gogf/gf/text/gregex"
)

// Decode parses <content> into and returns as map.
func Decode(content []byte) (map[string]interface{}, error) {
	res, err := convert(content)
	if err != nil {
		return nil, err
	}
	return mxj.NewMapXml(res)
}

// DecodeWithoutRoot parses <content> into a map, and returns the map without root level.
func DecodeWithoutRoot(content []byte) (map[string]interface{}, error) {
	res, err := convert(content)
	if err != nil {
		return nil, err
	}
	m, err := mxj.NewMapXml(res)
	if err != nil {
		return nil, err
	}
	for _, v := range m {
		if r, ok := v.(map[string]interface{}); ok {
			return r, nil
		}
	}
	return m, nil
}

// Encode encodes map <m> to a XML format content as bytes.
// The optional parameter <rootTag> is used to specify the XML root tag.
func Encode(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	return mxj.Map(m).Xml(rootTag...)
}

// Encode encodes map <m> to a XML format content as bytes with indent.
// The optional parameter <rootTag> is used to specify the XML root tag.
func EncodeWithIndent(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	return mxj.Map(m).XmlIndent("", "\t", rootTag...)
}

// ToJson converts <content> as XML format into JSON format bytes.
func ToJson(content []byte) ([]byte, error) {
	res, err := convert(content)
	if err != nil {
		return nil, err
	}
	mv, err := mxj.NewMapXml(res)
	if err == nil {
		return mv.Json()
	} else {
		return nil, err
	}
}

// convert converts the encoding of given XML content from XML root tag into UTF-8 encoding content.
func convert(xml []byte) (res []byte, err error) {
	patten := `<\?xml.*encoding\s*=\s*['|"](.*?)['|"].*\?>`
	matchStr, err := gregex.MatchString(patten, string(xml))
	if err != nil {
		return nil, err
	}
	xmlEncode := "UTF-8"
	if len(matchStr) == 2 {
		xmlEncode = matchStr[1]
	}
	xmlEncode = strings.ToUpper(xmlEncode)
	res, err = gregex.Replace(patten, []byte(""), xml)
	if err != nil {
		return nil, err
	}
	if xmlEncode != "UTF-8" && xmlEncode != "UTF8" {
		dst, err := gcharset.Convert("UTF-8", xmlEncode, string(res))
		if err != nil {
			return nil, err
		}
		res = []byte(dst)
	}
	return res, nil
}
