// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 数据校验.
// 本来打算取名gvalidator的，名字太长了，缩写一下
package gvalid

import (
    "strings"
    "regexp"
    "strconv"
    "gitee.com/johng/gf/third/github.com/fatih/structs"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/net/gipv4"
    "gitee.com/johng/gf/g/net/gipv6"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/container/gmap"
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
float                格式：float                                 说明：浮点数
boolean              格式：boolean                               说明：布尔值(1,true,on,yes:true | 0,false,off,no,"":false)
same                 格式：same:field                            说明：参数值必需与field参数的值相同
different            格式：different:field                       说明：参数值不能与field参数的值相同
in                   格式：in:value1,value2,...                  说明：参数值应该在value1,value2,...中(字符串匹配)
not-in               格式：not-in:value1,value2,...              说明：参数值不应该在value1,value2,...中(字符串匹配)
regex                格式：regex:pattern                         说明：参数值应当满足正则匹配规则pattern
*/

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

const (
    gSINGLE_RULE_PATTERN = `^([\w-]+):{0,1}(.*)` // 单条规则匹配正则
)

var (
    // 默认错误消息管理对象(并发安全)
    errorMsgMap  = gmap.NewStringStringMap()

    // 单规则正则对象，这里使用包内部变量存储，不需要多次解析
    ruleRegex, _ = regexp.Compile(gSINGLE_RULE_PATTERN)

    // 即时参数为空(nil|"")也需要校验的规则，主要是必需规则及关联规则
    mustCheckRulesEvenValueEmpty = map[string]struct{}{
        "required"             : struct{}{},
        "required-if"          : struct{}{},
        "required-unless"      : struct{}{},
        "required-with"        : struct{}{},
        "required-with-all"    : struct{}{},
        "required-without"     : struct{}{},
        "required-without-all" : struct{}{},
        "same"                 : struct{}{},
        "different"            : struct{}{},
        "in"                   : struct{}{},
        "not-in"               : struct{}{},
        "regex"                : struct{}{},
    }
)

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

// 判断必须字段
func checkRequired(value, ruleKey, ruleVal string, params map[string]string) bool {
    required := false
    switch ruleKey {
        // 必须字段
        case "required":
            required = true

        // 必须字段(当任意所给定字段值与所给值相等时)
        case "required-if":
            required = false
            array   := strings.Split(ruleVal, ",")
            // 必须为偶数，才能是键值对匹配
            if len(array)%2 == 0 {
                for i := 0; i < len(array); {
                    tk := array[i]
                    tv := array[i+1]
                    if v, ok := params[tk]; ok {
                        if strings.Compare(tv, v) == 0 {
                            required = true
                            break
                        }
                    }
                    i += 2
                }
            }

        // 必须字段(当所给定字段值与所给值都不相等时)
        case "required-unless":
            required = true
            array   := strings.Split(ruleVal, ",")
            // 必须为偶数，才能是键值对匹配
            if len(array)%2 == 0 {
                for i := 0; i < len(array); {
                    tk := array[i]
                    tv := array[i+1]
                    if v, ok := params[tk]; ok {
                        if strings.Compare(tv, v) == 0 {
                            required = false
                            break
                        }
                    }
                    i += 2
                }
            }

        // 必须字段(当所给定任意字段值不为空时)
        case "required-with":
            required = false
            array   := strings.Split(ruleVal, ",")
            for i := 0; i < len(array); i++ {
                if v, ok := params[array[i]]; ok {
                    if v != "" {
                        required = true
                        break
                    }
                }
            }

        // 必须字段(当所给定所有字段值都不为空时)
        case "required-with-all":
            required = true
            array   := strings.Split(ruleVal, ",")
            for i := 0; i < len(array); i++ {
                if v, ok := params[array[i]]; ok {
                    if v == "" {
                        required = false
                        break
                    }
                }
            }

        // 必须字段(当所给定任意字段值为空时)
        case "required-without":
            required = false
            array   := strings.Split(ruleVal, ",")
            for i := 0; i < len(array); i++ {
                if v, ok := params[array[i]]; ok {
                    if v == "" {
                        required = true
                        break
                    }
                }
            }

        // 必须字段(当所给定所有字段值都为空时)
        case "required-without-all":
            required = true
            array   := strings.Split(ruleVal, ",")
            for i := 0; i < len(array); i++ {
                if v, ok := params[array[i]]; ok {
                    if v != "" {
                        required = false
                        break
                    }
                }
            }
    }
    if required {
        return !(value == "")
    } else {
        return true
    }
}

// 对字段值长度进行检测
func checkLength(value, ruleKey, ruleVal string, custonMsgs map[string]string) string {
    msg := ""
    switch ruleKey {
        // 长度范围
        case "length":
            array := strings.Split(ruleVal, ",")
            min   := 0
            max   := 0
            if len(array) > 0 {
                if v, err := strconv.Atoi(strings.TrimSpace(array[0])); err == nil {
                    min = v
                }
            }
            if len(array) > 1 {
                if v, err := strconv.Atoi(strings.TrimSpace(array[1])); err == nil {
                    max = v
                }
            }
            if len(value) < min || len(value) > max {
                if v, ok := custonMsgs[ruleKey]; !ok {
                    msg = errorMsgMap.Get(ruleKey)
                } else {
                    msg = v
                }
                msg = strings.Replace(msg, ":min", strconv.Itoa(min), -1)
                msg = strings.Replace(msg, ":max", strconv.Itoa(max), -1)
                return msg
            }

        // 最小长度
        case "min-length":
            if min, err := strconv.Atoi(ruleVal); err == nil {
                if len(value) < min {
                    if v, ok := custonMsgs[ruleKey]; !ok {
                        msg = errorMsgMap.Get(ruleKey)
                    } else {
                        msg = v
                    }
                    msg = strings.Replace(msg, ":min", strconv.Itoa(min), -1)
                }
            } else {
                msg = "校验参数[" + ruleVal + "]应当为整数类型"
            }

        // 最大长度
        case "max-length":
            if max, err := strconv.Atoi(ruleVal); err == nil {
                if len(value) > max {
                    if v, ok := custonMsgs[ruleKey]; !ok {
                        msg = errorMsgMap.Get(ruleKey)
                    } else {
                        msg = v
                    }
                    msg = strings.Replace(msg, ":max", strconv.Itoa(max), -1)
                }
            } else {
                msg = "校验参数[" + ruleVal + "]应当为整数类型"
            }
    }
    return msg
}

// 对字段值大小进行检测
func checkSize(value, ruleKey, ruleVal string, custonMsgs map[string]string) string {
    msg := ""
    switch ruleKey {
        // 大小范围
        case "between":
            array := strings.Split(ruleVal, ",")
            min   := float64(0)
            max   := float64(0)
            if len(array) > 0 {
                if v, err := strconv.ParseFloat(strings.TrimSpace(array[0]), 10); err == nil {
                    min = v
                }
            }
            if len(array) > 1 {
                if v, err := strconv.ParseFloat(strings.TrimSpace(array[1]), 10); err == nil {
                    max = v
                }
            }
            if v, err := strconv.ParseFloat(value, 10); err == nil {
                if v < min || v > max {
                    if v, ok := custonMsgs[ruleKey]; !ok {
                        msg = errorMsgMap.Get(ruleKey)
                    } else {
                        msg = v
                    }
                    msg = strings.Replace(msg, ":min", strconv.FormatFloat(min, 'f', -1, 64), -1)
                    msg = strings.Replace(msg, ":max", strconv.FormatFloat(max, 'f', -1, 64), -1)
                }
            } else {
                msg = "输入参数[" + value + "]应当为数字类型"
            }

        // 最小值
        case "min":
            if min, err := strconv.ParseFloat(ruleVal, 10); err == nil {
                if v, err := strconv.ParseFloat(value, 10); err == nil {
                    if v < min {
                        if v, ok := custonMsgs[ruleKey]; !ok {
                            msg = errorMsgMap.Get(ruleKey)
                        } else {
                            msg = v
                        }
                        msg = strings.Replace(msg, ":min", strconv.FormatFloat(min, 'f', -1, 64), -1)
                    }
                } else {
                    msg = "输入参数[" + value + "]应当为数字类型"
                }
            } else {
                msg = "校验参数[" + ruleVal + "]应当为数字类型"
            }

        // 最大值
        case "max":
            if max, err := strconv.ParseFloat(ruleVal, 10); err == nil {
                if v, err := strconv.ParseFloat(value, 10); err == nil {
                    if v > max {
                        if v, ok := custonMsgs[ruleKey]; !ok {
                            msg = errorMsgMap.Get(ruleKey)
                        } else {
                            msg = v
                        }
                        msg = strings.Replace(msg, ":max", strconv.FormatFloat(max, 'f', -1, 64), -1)
                    }
                } else {
                    msg = "输入参数[" + value + "]应当为数字类型"
                }
            } else {
                msg = "校验参数[" + ruleVal + "]应当为数字类型"
            }
    }
    return msg
}

// 检测键值对参数Map，注意返回参数是一个2维的关联数组，第一维键名为参数键名，第二维为带有错误的校验规则名称，值为错误信息
func CheckMap(params map[string]interface{}, rules map[string]string, msgs...map[string]interface{}) Error {
    var value interface{}
    // 自定义消息，非必须参数，因此这里需要做判断
    customMsgs := make(map[string]interface{})
    if len(msgs) > 0 && len(msgs[0]) > 0 {
        customMsgs = msgs[0]
    }
    errorMsgs := make(map[string]map[string]string)
    // 以校验规则作为基础
    for key, rule := range rules {
        value = nil
        if v, ok := params[key]; ok {
            value = v
        }
        msg, _ := customMsgs[key]
        if e := Check(value, rule, msg, params); e != nil {
            _, m := e.FirstItem()
            // 如果值为nil|""，并且不需要require*验证时，其他验证失效
            if value == nil || gconv.String(value) == "" {
                required := false
                // rule => error
                for k, _ := range m {
                    if _, ok := mustCheckRulesEvenValueEmpty[k]; ok {
                        required = true
                        break
                    }
                }
                if !required {
                    continue
                }
            }
            if _, ok := errorMsgs[key]; !ok {
                errorMsgs[key] = make(map[string]string)
            }
            for k, v := range m {
                errorMsgs[key][k] = v
            }
        }
    }
    if len(errorMsgs) > 0 {
        return errorMsgs
    }
    return nil
}

// 校验struct对象属性，object参数也可以是一个指向对象的指针，返回值同CheckMap方法
func CheckStruct(st interface{}, rules map[string]string, msgs...map[string]interface{}) Error {
    fields := structs.Fields(st)
    if rules == nil {
        rules = make(map[string]string)
    }
    params  := make(map[string]interface{})
    errMsgs := (map[string]interface{})(nil)
    if len(msgs) == 0 {
        errMsgs = make(map[string]interface{})
    } else {
        errMsgs = msgs[0]
    }
    for _, field := range fields {
        params[field.Name()] = field.Value()
        if tag := field.Tag("gvalid"); tag != "" {
            match, _ := gregex.MatchString(`\s*((\w+)\s*@){0,1}\s*([^#]+)\s*(#\s*(.*)){0,1}\s*`, tag)
            name := match[2]
            rule := match[3]
            msg  := match[5]
            if len(name) == 0 {
                name = field.Name()
            }
            // params参数使用别名**扩容**(而不仅仅使用别名)，仅用于验证使用
            if _, ok := params[name]; !ok {
                params[name] = field.Value()
            }
            // 校验规则
            if _, ok := rules[name]; !ok {
                rules[name] = rule
            }
            // 错误提示
            if len(msg) > 0 {
                ruleArray := strings.Split(rule, "|")
                msgArray  := strings.Split(msg, "|")
                for k, v := range ruleArray {
                    if len(msgArray[k]) == 0 {
                        continue
                    }
                    array := strings.Split(v, ":")
                    if _, ok := errMsgs[name]; !ok {
                        errMsgs[name] = make(map[string]string)
                    }
                    errMsgs[name].(map[string]string)[array[0]] = msgArray[k]
                }
            }
        }

    }
    return CheckMap(params, rules, errMsgs)
}

// 检测单条数据的规则.
// val为校验数据，可以为任意基本数据格式；
// msgs为自定义错误信息，由于同一条数据的校验规则可能存在多条，为方便调用，参数类型支持string/map[string]string，允许传递多个自定义的错误信息，如果类型为string，那么中间使用"|"符号分隔多个自定义错误；
// params参数为表单联合校验参数，对于需要联合校验的规则有效，如：required-*、same、different；
func Check(val interface{}, rules string, msgs interface{}, params...map[string]interface{}) Error {
    // 内部会将参数全部转换为字符串类型进行校验
    value     := strings.TrimSpace(gconv.String(val))
    data      := make(map[string]string)
    errorMsgs := make(map[string]string)
    if len(params) > 0 {
        for k, v := range params[0] {
            data[k] = gconv.String(v)
        }
    }
    // 自定义错误消息处理
    list       := make([]string, 0)
    custonMsgs := make(map[string]string)
    switch value := msgs.(type) {
        case map[string]string:
            custonMsgs = value
        case string:
            list = strings.Split(value, "|")
    }
    items := strings.Split(strings.TrimSpace(rules), "|")
    for index := 0; index < len(items); {
        item    := items[index]
        results := ruleRegex.FindStringSubmatch(item)
        ruleKey := strings.TrimSpace(results[1])
        ruleVal := strings.TrimSpace(results[2])
        match   := false
        if len(list) > index {
            custonMsgs[ruleKey] = strings.TrimSpace(list[index])
        }
        switch ruleKey {
            // 必须字段
            case "required":          fallthrough
            case "required-if":       fallthrough
            case "required-unless":   fallthrough
            case "required-with":     fallthrough
            case "required-with-all": fallthrough
            case "required-without":  fallthrough
            case "required-without-all":
                match = checkRequired(value, ruleKey, ruleVal, data)

            // 长度范围
            case "length":            fallthrough
            case "min-length":        fallthrough
            case "max-length":
                if msg := checkLength(value, ruleKey, ruleVal, custonMsgs); msg != "" {
                    errorMsgs[ruleKey] = msg
                } else {
                    match = true
                }

            // 大小范围
            case "min":               fallthrough
            case "max":               fallthrough
            case "between":
                if msg := checkSize(value, ruleKey, ruleVal, custonMsgs); msg != "" {
                    errorMsgs[ruleKey] = msg
                } else {
                    match = true
                }

            // 自定义正则判断
            case "regex":
                // 需要判断是否被|符号截断，如果是，那么需要进行整合
                for i := index + 1; i < len(items); i++ {
                    // 判断下一个规则是否合法，不合法那么和当前正则规则进行整合
                    if !gregex.IsMatchString(gSINGLE_RULE_PATTERN, items[i]) {
                        ruleVal += "|" + items[i]
                        index++
                    }
                }
                match = gregex.IsMatchString(ruleVal, value)

            // 日期格式，
            case "date":
                if _, err := gtime.StrToTime(value); err == nil {
                    match = true
                    break
                }

            // 日期格式，需要给定日期格式
            case "date-format":
                if _, err := gtime.StrToTimeFormat(value, ruleVal); err == nil {
                    match = true
                }

            // 两字段值应相同(非敏感字符判断，非类型判断)
            case "same":
                if v, ok := data[ruleVal]; ok {
                    if strings.Compare(value, v) == 0 {
                        match = true
                    }
                }

            // 两字段值不应相同(非敏感字符判断，非类型判断)
            case "different":
                match = true
                if v, ok := data[ruleVal]; ok {
                    if strings.Compare(value, v) == 0 {
                        match = false
                    }
                }

            // 字段值应当在指定范围中
            case "in":
                array := strings.Split(ruleVal, ",")
                for _, v := range array {
                    if strings.Compare(value, strings.TrimSpace(v)) == 0 {
                        match = true
                        break
                    }
                }

            // 字段值不应当在指定范围中
            case "not-in":
                match  = true
                array := strings.Split(ruleVal, ",")
                for _, v := range array {
                    if strings.Compare(value, strings.TrimSpace(v)) == 0 {
                        match = false
                        break
                    }
                }

            /*
             * 验证所给手机号码是否符合手机号的格式.
             * 移动：134、135、136、137、138、139、150、151、152、157、158、159、182、183、184、187、188、178(4G)、147(上网卡)；
             * 联通：130、131、132、155、156、185、186、176(4G)、145(上网卡)、175；
             * 电信：133、153、180、181、189 、177(4G)；
             * 卫星通信：  1349
             * 虚拟运营商：170、173
             */
            case "phone":
                match = gregex.IsMatchString(`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^17[0,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$`, value)

            // 国内座机电话号码："XXXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"、"XXXXXXXX"
            case "telephone":
                match = gregex.IsMatchString(`^((\d{3,4})|\d{3,4}-)?\d{7,8}$`, value)

            // 腾讯QQ号，从10000开始
            case "qq":
                match = gregex.IsMatchString(`^[1-9][0-9]{4,}$`, value)

                // 中国邮政编码
            case "postcode":
                match = gregex.IsMatchString(`^[1-9]\d{5}$`, value)

            /*
                公民身份证号
                xxxxxx yyyy MM dd 375 0     十八位
                xxxxxx   yy MM dd  75 0     十五位

                地区：[1-9]\d{5}
                年的前两位：(18|19|([23]\d))      1800-2399
                年的后两位：\d{2}
                月份：((0[1-9])|(10|11|12))
                天数：(([0-2][1-9])|10|20|30|31) 闰年不能禁止29+

                三位顺序码：\d{3}
                两位顺序码：\d{2}
                校验码：   [0-9Xx]

                十八位：^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$
                十五位：^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{2}$

                总：
                (^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{2}$)
             */
            case "id-number":
                match = gregex.IsMatchString(`(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{2}$)`, value)

            // 通用帐号规则(字母开头，只能包含字母、数字和下划线，长度在6~18之间)
            case "passport":
                match = gregex.IsMatchString(`^[a-zA-Z]{1}\w{5,17}$`, value)

            // 通用密码(任意可见字符，长度在6~18之间)
            case "password":
                match = gregex.IsMatchString(`^[\w\S]{6,18}$`, value)

            // 中等强度密码(在弱密码的基础上，必须包含大小写字母和数字)
            case "password2":
                if gregex.IsMatchString(`^[\w\S]{6,18}$`, value)  && gregex.IsMatchString(`[a-z]+`, value) && gregex.IsMatchString(`[A-Z]+`, value) && gregex.IsMatchString(`\d+`, value) {
                    match = true
                }

            // 强等强度密码(在弱密码的基础上，必须包含大小写字母、数字和特殊字符)
            case "password3":
                if gregex.IsMatchString(`^[\w\S]{6,18}$`, value) && gregex.IsMatchString(`[a-z]+`, value) && gregex.IsMatchString(`[A-Z]+`, value) && gregex.IsMatchString(`\d+`, value) && gregex.IsMatchString(`\S+`, value) {
                    match = true
                }

            // json
            case "json":
                if _, err := gjson.Decode([]byte(value)); err == nil {
                    match = true
                }

            // 整数
            case "integer":
                if _, err := strconv.Atoi(value); err == nil {
                    match = true
                }

            // 小数
            case "float":
                if _, err := strconv.ParseFloat(value, 10); err == nil {
                    match = true
                }

            // 布尔值(1,true,on,yes:true | 0,false,off,no,"":false)
            case "boolean":
                if value != "" && value != "0" && value != "false" && value != "off" && value != "no" {
                    match = true
                }

            // 邮件
            case "email":
                match = gregex.IsMatchString(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, value)

            // URL
            case "url":
                match = gregex.IsMatchString(`^(?:([A-Za-z]+):)?(\/{0,3})([0-9.\-A-Za-z]+)(?::(\d+))?(?:\/([^?#]*))?(?:\?([^#]*))?(?:#(.*))?$`, value)

            // domain
            case "domain":
                match = gregex.IsMatchString(`^([0-9a-zA-Z][0-9a-zA-Z-]{0,62}\.)+([0-9a-zA-Z][0-9a-zA-Z-]{0,62})\.?$`, value)

            // IP(IPv4/IPv6)
            case "ip":
                match = gipv4.Validate(value) || gipv6.Validate(value)

            // IPv4
            case "ipv4":
                match = gipv4.Validate(value)

            // IPv6
            case "ipv6":
                match = gipv6.Validate(value)

            // MAC地址
            case "mac":
                match = gregex.IsMatchString(`^([0-9A-Fa-f]{2}-){5}[0-9A-Fa-f]{2}$`, value)

            default:
                errorMsgs[ruleKey] = "Invalid rule name:" + ruleKey
        }

        // 错误消息整合
        if !match {
            // 不存在则使用默认的错误信息，
            // 如果在校验过程中已经设置了错误信息，那么这里便不作处理
            if _, ok := errorMsgs[ruleKey]; !ok {
                if msg, ok := custonMsgs[ruleKey]; ok {
                    errorMsgs[ruleKey] = msg
                } else {
                    errorMsgs[ruleKey] = errorMsgMap.Get(ruleKey)
                }
            }
        }
        index++
    }
    if len(errorMsgs) > 0 {
        e := make(Error)
        e[value] = errorMsgs
        return e
    }
    return nil
}