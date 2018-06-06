// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// XML
package gxml

import (
    "github.com/clbanning/mxj"
    "encoding/xml"
    "io"
	"gitee.com/johng/gf/g/util/gregx"
	"github.com/axgle/mahonia"
	"errors"
	"fmt"
)

// 将XML内容解析为map变量
func Decode(xmlbyte []byte) (map[string]interface{}, error) {
    Prepare(xmlbyte)
    return mxj.NewMapXml(xmlbyte)
}

// 将map变量解析为XML格式内容
func Encode(v map[string]interface{}, rootTag...string) ([]byte, error) {
    return mxj.Map(v).Xml(rootTag...)
}

func EncodeWithIndent(v map[string]interface{}, rootTag...string) ([]byte, error) {
    return mxj.Map(v).XmlIndent("", "\t", rootTag...)
}

// XML格式内容直接转换为JSON格式内容
func ToJson(xmlbyte []byte) ([]byte, error) {
    Prepare(xmlbyte)
	mv, err := mxj.NewMapXml(xmlbyte)
	if err == nil {
        return mv.Json()
    } else {
        return nil, err
    }
}

//XML字符集预处理
//@author wenzi1 
//@date 20180604
func Prepare(xmlbyte []byte) error {
	patten := "<\\?xml\\s+version\\s*=.*?\\s+encoding\\s*=\\s*[\\'|\"](.*?)[\\'|\"]\\s*\\?\\s*>"
	charsetReader := func(charset string, input io.Reader) (io.Reader, error) {
		reader := mahonia.GetCharset(charset)
		if reader == nil {
			return nil, errors.New(fmt.Sprintf("not support charset:%s", charset))
		}
		return reader.NewDecoder().NewReader(input), nil
	}

	matchStr, err := gregx.MatchString(patten, string(xmlbyte))
	if err != nil {
		return err
	}

	charset := mahonia.GetCharset(matchStr[1])
	if charset == nil {
		return errors.New(fmt.Sprintf("not support charset:%s", matchStr[1]))
	}

	if charset.Name != "UTF-8" {
		mxj.CustomDecoder = &xml.Decoder{Strict:false,CharsetReader:charsetReader}
	}
	return nil
}