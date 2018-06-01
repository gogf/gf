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
    "strings"
    "golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 将XML内容解析为map变量
func Decode(xmlbyte []byte) (map[string]interface{}, error) {
	if strings.Index(string(xmlbyte), "encoding=\"UTF-8\"") == -1 {
		charsetReader := func(charset string, input io.Reader) (io.Reader, error) {
			reader := input
			switch charset {
			case "GBK":
				reader = transform.NewReader(input,simplifiedchinese.GBK.NewDecoder())
			case "GB18030":
				reader = transform.NewReader(input, simplifiedchinese.GB18030.NewDecoder())
			default:
				reader = input
			}
			return reader, nil
		}	
		mxj.CustomDecoder = &xml.Decoder{Strict:false,CharsetReader:charsetReader}
	}
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
	//input, _ := ioutil.ReadAll(bytes.NewReader(xmlbyte))
	//input:= bytes.NewReader(xmlbyte)
	//charset := "UTF-8"
	//@wenzi1 20180529
	//XML中的encoding如果非UTF-8时mxj包中会自动转换，但是自动转换时需要调用方提供转换的方法,这里提供一个空方法
	//所以如果涉及到字符集转换那么需要用户自行转为utf8时再调用该方法
	if strings.Index(string(xmlbyte), "encoding=\"UTF-8\"") == -1 { 
		charsetReader := func(charset string, input io.Reader) (io.Reader, error) {
			reader := input
			switch charset {
			case "GBK":
				reader = transform.NewReader(input,simplifiedchinese.GBK.NewDecoder())
			case "GB18030":
				reader = transform.NewReader(input, simplifiedchinese.GB18030.NewDecoder())
			default:
				reader = input
			}
			return reader, nil
		}	
		mxj.CustomDecoder = &xml.Decoder{Strict:false,CharsetReader:charsetReader}
	}
	mv, err := mxj.NewMapXml(xmlbyte)
	
	if err == nil {
        return mv.Json()
    } else {
        return nil, err
    }
}