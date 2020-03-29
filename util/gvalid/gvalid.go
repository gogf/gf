// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gvalid implements powerful and useful data/form validation functionality.
package gvalid

import (
	"strings"

	"github.com/gogf/gf/text/gregex"
)

/*
参考：https://laravel.com/docs/5.5/validation#available-validation-rules
规则如下：
required             格式：required                              说明：必需参数
required-if          格式：required-if:field,value,...           说明：必需参数(当任意所给定字段值与所给值相等时，即：当field字段的值为value时，当前验证字段为必须参数)
required-unless      格式：required-unless:field,value,...       说明：必需参数(当所给定字段值与所给值都不相等时，即：当field字段的值不为value时，当前验证字段为必须参数)
required-with        格式：required-with:field1,field2,...       说明：必需参数(当所给定任意字段值不为空时)
required-with-all    格式：required-with-all:field1,field2,...   说明：必须参数(当所给定所有字段值都不为空时)
required-without     格式：required-without:field1,field2,...    说明：必需参数(当所给定任意字段值为空时)
required-without-all 格式：required-without-all:field1,field2,...说明：必须参数(当所给定所有字段值都为空时)
date                 格式：date                                  说明：参数为常用日期类型，格式：2006-01-02, 20060102, 2006.01.02
date-format          格式：date-format:format                    说明：判断日期是否为指定的日期格式，format为Go日期格式(可以包含时间)
email                格式：email                                 说明：EMAIL邮箱地址
phone                格式：phone                                 说明：手机号
telephone            格式：telephone                             说明：国内座机电话号码，"XXXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"、"XXXXXXXX"
passport             格式：passport                              说明：通用帐号规则(字母开头，只能包含字母、数字和下划线，长度在6~18之间)
password             格式：password                              说明：通用密码(任意可见字符，长度在6~18之间)
password2            格式：password2                             说明：中等强度密码(在弱密码的基础上，必须包含大小写字母和数字)
password3            格式：password3                             说明：强等强度密码(在弱密码的基础上，必须包含大小写字母、数字和特殊字符)
postcode             格式：postcode                              说明：中国邮政编码
id-number            格式：id-number                             说明：公民身份证号码
luhn                 格式：luhn                                  说明：银行号验证
qq                   格式：qq                                    说明：腾讯QQ号码
ip                   格式：ip                                    说明：IPv4/IPv6地址
ipv4                 格式：ipv4                                  说明：IPv4地址
ipv6                 格式：ipv6                                  说明：IPv6地址
mac                  格式：mac                                   说明：MAC地址
url                  格式：url                                   说明：URL
domain               格式：domain                                说明：域名
length               格式：length:min,max                        说明：参数长度为min到max(长度参数为整形)，注意中文一个汉字占3字节
min-length           格式：min-length:min                        说明：参数长度最小为min(长度参数为整形)，注意中文一个汉字占3字节
max-length           格式：max-length:max                        说明：参数长度最大为max(长度参数为整形)，注意中文一个汉字占3字节
between              格式：between:min,max                       说明：参数大小为min到max(支持整形和浮点类型参数)
min                  格式：min:min                               说明：参数最小为min(支持整形和浮点类型参数)
max                  格式：max:max                               说明：参数最大为max(支持整形和浮点类型参数)
json                 格式：json                                  说明：判断数据格式为JSON
integer              格式：integer                               说明：整数
float                格式：float                                 说明：浮点数(整数也是浮点数)
boolean              格式：boolean                               说明：布尔值(1,true,on,yes:true | 0,false,off,no,"":false)
same                 格式：same:field                            说明：参数值必需与field参数的值相同
different            格式：different:field                       说明：参数值不能与field参数的值相同
in                   格式：in:value1,value2,...                  说明：参数值应该在value1,value2,...中(字符串匹配)
not-in               格式：not-in:value1,value2,...              说明：参数值不应该在value1,value2,...中(字符串匹配)
regex                格式：regex:pattern                         说明：参数值应当满足正则匹配规则pattern
*/

// 自定义错误信息: map[键名] => 字符串|map[规则]错误信息
type CustomMsg = map[string]interface{}

// FieldRangeMsg 用于field-in,field-not-in的数据信息: map[键名] => []string range部分
type FieldRangeMsg = map[string][]string

// 解析单条sequence tag，格式: [alias@]rule[...#msg...]
func parseSequenceTag(tag string) (name, rule, msg string) {
	// Complete sequence tag.
	// Eg: required|length:2,20|password3|same:password1#||密码强度不足|两次密码不一致
	match, _ := gregex.MatchString(`\s*((\w+)\s*@){0,1}\s*([^#]+)\s*(#\s*(.*)){0,1}\s*`, tag)
	return strings.TrimSpace(match[2]), strings.TrimSpace(match[3]), strings.TrimSpace(match[5])
}
