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
)

// 将XML内容解析为map变量,入口参数使用的字符集必须为UTF-8的
func Decode(xmlbyte []byte) (map[string]interface{}, error) {
    //@author wenzi1  
    //@date   2018-05-29
	//XML中的encoding如果非UTF-8时mxj包中会自动转换，但是自动转换时需要调用方提供转换的方法,这里提供一个空方法
	//所以如果涉及到字符集转换那么需要用户自行转为utf8时再调用该方法
    charsetReader := func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}	
	mxj.CustomDecoder = &xml.Decoder{CharsetReader:charsetReader}
    return mxj.NewMapXml(xmlbyte)
}

// 将map变量解析为XML格式内容
func Encode(v map[string]interface{}, rootTag...string) ([]byte, error) {
    return mxj.Map(v).Xml(rootTag...)
}

func EncodeWithIndent(v map[string]interface{}, rootTag...string) ([]byte, error) {
    return mxj.Map(v).XmlIndent("", "\t", rootTag...)
}

// XML格式内容直接转换为JSON格式内容,入口参数使用的字符集必须为UTF-8的
func ToJson(xmlbyte []byte) ([]byte, error) {
    //@author wenzi1  
    //@date   2018-05-29
	//XML中的encoding如果非UTF-8时mxj包中会自动转换，但是自动转换时需要调用方提供转换的方法,这里提供一个空方法
	//所以如果涉及到字符集转换那么需要用户自行转为utf8时再调用该方法
	charsetReader := func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}	
	mxj.CustomDecoder = &xml.Decoder{CharsetReader:charsetReader}
    if mv, err := mxj.NewMapXml(xmlbyte); err == nil {
        return mv.Json()
    } else {
        return nil, err
    }
}