// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 默认的错误消息定义。

package gvalid

// 默认规则校验错误消息(可以通过接口自定义错误消息)
var defaultMessages = map[string]string {
    "required"             : "字段不能为空",
    "required-if"          : "字段不能为空",
    "required-unless"      : "字段不能为空",
    "required-with"        : "字段不能为空",
    "required-with-all"    : "字段不能为空",
    "required-without"     : "字段不能为空",
    "required-without-all" : "字段不能为空",
    "date"                 : "日期格式不正确",
    "date-format"          : "日期格式不正确",
    "email"                : "邮箱地址格式不正确",
    "phone"                : "手机号码格式不正确",
    "telephone"            : "电话号码格式不正确",
    "passport"             : "账号格式不合法，必需以字母开头，只能包含字母、数字和下划线，长度在6~18之间",
    "password"             : "密码格式不合法，密码格式为任意6-18位的可见字符",
    "password2"            : "密码格式不合法，密码格式为任意6-18位的可见字符，必须包含大小写字母和数字",
    "password3"            : "密码格式不合法，密码格式为任意6-18位的可见字符，必须包含大小写字母、数字和特殊字符",
    "postcode"             : "邮政编码不正确",
    "id-number"            : "身份证号码不正确",
    "qq"                   : "QQ号码格式不正确",
    "ip"                   : "IP地址格式不正确",
    "ipv4"                 : "IPv4地址格式不正确",
    "ipv6"                 : "IPv6地址格式不正确",
    "mac"                  : "MAC地址格式不正确",
    "url"                  : "URL地址格式不正确",
    "domain"               : "域名格式不正确",
    "length"               : "字段长度为:min到:max个字符",
    "min-length"           : "字段最小长度为:min",
    "max-length"           : "字段最大长度为:max",
    "between"              : "字段大小为:min到:max",
    "min"                  : "字段最小值为:min",
    "max"                  : "字段最大值为:max",
    "json"                 : "字段应当为JSON格式",
    "xml"                  : "字段应当为XML格式",
    "array"                : "字段应当为数组",
    "integer"              : "字段应当为整数",
    "float"                : "字段应当为浮点数",
    "boolean"              : "字段应当为布尔值",
    "same"                 : "字段值不合法",
    "different"            : "字段值不合法",
    "in"                   : "字段值不合法",
    "not-in"               : "字段值不合法",
    "regex"                : "字段值不合法",
}

// 初始化错误消息管理对象
func init() {
    errorMsgMap.BatchSet(defaultMessages)
}

// 替换默认的错误提示为指定的自定义提示
// 主要作用：
// 1、便于多语言错误提示设置；
// 2、默认错误提示信息不满意；
func SetDefaultErrorMsgs(msgs map[string]string) {
    errorMsgMap.BatchSet(msgs)
}