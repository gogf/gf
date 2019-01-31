// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gvalid

import (
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/net/gipv4"
    "gitee.com/johng/gf/g/net/gipv6"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gregex"
    "regexp"
    "strconv"
    "strings"
)

const (
    gSINGLE_RULE_PATTERN = `^([\w-]+):{0,1}(.*)` // 单条规则匹配正则
)

var (
    // 默认错误消息管理对象(并发安全)
    errorMsgMap  = gmap.NewStringStringMap()

    // 单规则正则对象，这里使用包内部变量存储，不需要多次解析
    ruleRegex, _ = regexp.Compile(gSINGLE_RULE_PATTERN)

    // 即使参数为空(nil|"")也需要校验的规则，主要是必需规则及关联规则
    mustCheckRulesEvenValueEmpty = map[string]struct{} {
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
    // 所有支持的校验规则
    allSupportedRules = map[string]struct{} {
        "required"                  : struct{}{},
        "required-if"               : struct{}{},
        "required-unless"           : struct{}{},
        "required-with"             : struct{}{},
        "required-with-all"         : struct{}{},
        "required-without"          : struct{}{},
        "required-without-all"      : struct{}{},
        "date"                      : struct{}{},
        "date-format"               : struct{}{},
        "email"                     : struct{}{},
        "phone"                     : struct{}{},
        "telephone"                 : struct{}{},
        "passport"                  : struct{}{},
        "password"                  : struct{}{},
        "password2"                 : struct{}{},
        "password3"                 : struct{}{},
        "postcode"                  : struct{}{},
        "id-number"                 : struct{}{},
        "qq"                        : struct{}{},
        "ip"                        : struct{}{},
        "ipv4"                      : struct{}{},
        "ipv6"                      : struct{}{},
        "mac"                       : struct{}{},
        "url"                       : struct{}{},
        "domain"                    : struct{}{},
        "length"                    : struct{}{},
        "min-length"                : struct{}{},
        "max-length"                : struct{}{},
        "between"                   : struct{}{},
        "min"                       : struct{}{},
        "max"                       : struct{}{},
        "json"                      : struct{}{},
        "integer"                   : struct{}{},
        "float"                     : struct{}{},
        "boolean"                   : struct{}{},
        "same"                      : struct{}{},
        "different"                 : struct{}{},
        "in"                        : struct{}{},
        "not-in"                    : struct{}{},
        "regex"                     : struct{}{},
    }
    // 布尔Map
    boolMap = map[string]struct{} {
        // true
        "1"     : struct{}{},
        "true"  : struct{}{},
        "on"    : struct{}{},
        "yes"   : struct{}{},
        // false
        ""      : struct{}{},
        "0"     : struct{}{},
        "false" : struct{}{},
        "off"   : struct{}{},
        "no"    : struct{}{},
    }
)

// 检测单条数据的规则:
// value为需要校验的数据，可以为任意基本数据类型；
// msgs为自定义错误信息，由于同一条数据的校验规则可能存在多条，为方便调用，参数类型支持 string/map[string]string ，允许传递多个自定义的错误信息，如果类型为string，那么中间使用"|"符号分隔多个自定义错误；
// params参数为联合校验参数，对于需要联合校验的规则有效，如：required-*、same、different；
func Check(value interface{}, rules string, msgs interface{}, params...map[string]interface{}) *Error {
    // 内部会将参数全部转换为字符串类型进行校验
    val       := strings.TrimSpace(gconv.String(value))
    data      := make(map[string]string)
    errorMsgs := make(map[string]string)
    if len(params) > 0 {
        for k, v := range params[0] {
            data[k] = gconv.String(v)
        }
    }
    // 自定义错误消息处理
    msgArray     := make([]string, 0)
    customMsgMap := make(map[string]string)
    switch v := msgs.(type) {
        case map[string]string:
            customMsgMap = v

        case string:
            msgArray = strings.Split(v, "|")
    }
    ruleItems := strings.Split(strings.TrimSpace(rules), "|")
    // 规则项预处理, 主要解决规则中存在的"|"关键字符号
    for i := 0; ; {
        array := strings.Split(ruleItems[i], ":")
        if _, ok := allSupportedRules[array[0]]; !ok {
            if i > 0 {
                ruleItems[i - 1] += "|" + ruleItems[i]
                ruleItems = append(ruleItems[ : i], ruleItems[i + 1 : ]...)
            } else {
                return newErrorStr("invalid_rules", "invalid rules:" + rules)
            }
        } else {
            i++
        }
        if i == len(ruleItems) {
            break
        }
    }
    for index := 0; index < len(ruleItems); {
        item    := ruleItems[index]
        results := ruleRegex.FindStringSubmatch(item)
        ruleKey := strings.TrimSpace(results[1])
        ruleVal := strings.TrimSpace(results[2])
        match   := false
        if len(msgArray) > index {
            customMsgMap[ruleKey] = strings.TrimSpace(msgArray[index])
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
                match = checkRequired(val, ruleKey, ruleVal, data)

            // 长度范围
            case "length":            fallthrough
            case "min-length":        fallthrough
            case "max-length":
                if msg := checkLength(val, ruleKey, ruleVal, customMsgMap); msg != "" {
                    errorMsgs[ruleKey] = msg
                } else {
                    match = true
                }

            // 大小范围
            case "min":               fallthrough
            case "max":               fallthrough
            case "between":
                if msg := checkSize(val, ruleKey, ruleVal, customMsgMap); msg != "" {
                    errorMsgs[ruleKey] = msg
                } else {
                    match = true
                }

            // 自定义正则判断
            case "regex":
                // 需要判断是否被|符号截断，如果是，那么需要进行整合
                for i := index + 1; i < len(ruleItems); i++ {
                    // 判断下一个规则是否合法，不合法那么和当前正则规则进行整合
                    if !gregex.IsMatchString(gSINGLE_RULE_PATTERN, ruleItems[i]) {
                        ruleVal += "|" + ruleItems[i]
                        index++
                    }
                }
                match = gregex.IsMatchString(ruleVal, val)

            // 日期格式，
            case "date":
                // 使用标准日期格式检查，但是日期之间必须带连接符号
                if _, err := gtime.StrToTime(val); err == nil {
                    match = true
                    break
                }
                // 检查是否不带日期连接符号的格式
                if _, err := gtime.StrToTime(val, "Ymd"); err == nil {
                    match = true
                    break
                }

            // 日期格式，需要给定日期格式
            case "date-format":
                if _, err := gtime.StrToTimeFormat(val, ruleVal); err == nil {
                    match = true
                }

            // 两字段值应相同(非敏感字符判断，非类型判断)
            case "same":
                if v, ok := data[ruleVal]; ok {
                    if strings.Compare(val, v) == 0 {
                        match = true
                    }
                }

            // 两字段值不应相同(非敏感字符判断，非类型判断)
            case "different":
                match = true
                if v, ok := data[ruleVal]; ok {
                    if strings.Compare(val, v) == 0 {
                        match = false
                    }
                }

            // 字段值应当在指定范围中
            case "in":
                array := strings.Split(ruleVal, ",")
                for _, v := range array {
                    if strings.Compare(val, strings.TrimSpace(v)) == 0 {
                        match = true
                        break
                    }
                }

            // 字段值不应当在指定范围中
            case "not-in":
                match  = true
                array := strings.Split(ruleVal, ",")
                for _, v := range array {
                    if strings.Compare(val, strings.TrimSpace(v)) == 0 {
                        match = false
                        break
                    }
                }

            /*
             * 验证所给手机号码是否符合手机号的格式.
             * 移动: 134、135、136、137、138、139、150、151、152、157、158、159、182、183、184、187、188、178(4G)、147(上网卡)；
             * 联通: 130、131、132、155、156、185、186、176(4G)、145(上网卡)、175；
             * 电信: 133、153、180、181、189 、177(4G)；
             * 卫星通信:  1349
             * 虚拟运营商: 170、173
             * 2018新增: 16x, 19x
             */
            case "phone":
                match = gregex.IsMatchString(`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^16[\d]{9}$|^17[0,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$|^19[\d]{9}$`, val)

            // 国内座机电话号码："XXXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"、"XXXXXXXX"
            case "telephone":
                match = gregex.IsMatchString(`^((\d{3,4})|\d{3,4}-)?\d{7,8}$`, val)

            // 腾讯QQ号，从10000开始
            case "qq":
                match = gregex.IsMatchString(`^[1-9][0-9]{4,}$`, val)

            // 中国邮政编码
            case "postcode":
                match = gregex.IsMatchString(`^\d{6}$`, val)

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
                十五位：^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$

                总：
                (^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)
             */
            case "id-number":
                match = gregex.IsMatchString(`(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)`, val)

            // 通用帐号规则(字母开头，只能包含字母、数字和下划线，长度在6~18之间)
            case "passport":
                match = gregex.IsMatchString(`^[a-zA-Z]{1}\w{5,17}$`, val)

            // 通用密码(任意可见字符，长度在6~18之间)
            case "password":
                match = gregex.IsMatchString(`^[\w\S]{6,18}$`, val)

            // 中等强度密码(在弱密码的基础上，必须包含大小写字母和数字)
            case "password2":
                if gregex.IsMatchString(`^[\w\S]{6,18}$`, val)  && gregex.IsMatchString(`[a-z]+`, val) && gregex.IsMatchString(`[A-Z]+`, val) && gregex.IsMatchString(`\d+`, val) {
                    match = true
                }

            // 强等强度密码(在弱密码的基础上，必须包含大小写字母、数字和特殊字符)
            case "password3":
                if gregex.IsMatchString(`^[\w\S]{6,18}$`, val) && gregex.IsMatchString(`[a-z]+`, val) && gregex.IsMatchString(`[A-Z]+`, val) && gregex.IsMatchString(`\d+`, val) && gregex.IsMatchString(`[^a-zA-Z0-9]+`, val) {
                    match = true
                }

            // json
            case "json":
                if _, err := gjson.Decode([]byte(val)); err == nil {
                    match = true
                }

            // 整数
            case "integer":
                if _, err := strconv.Atoi(val); err == nil {
                    match = true
                }

            // 小数
            case "float":
                if _, err := strconv.ParseFloat(val, 10); err == nil {
                    match = true
                }

            // 布尔值(1,true,on,yes:true | 0,false,off,no,"":false)
            case "boolean":
                match = false
                if _, ok := boolMap[strings.ToLower(val)]; ok {
                    match = true
                }

            // 邮件
            case "email":
                match = gregex.IsMatchString(`^[a-zA-Z0-9_\-\.]+@[a-zA-Z0-9_\-]+(\.[a-zA-Z0-9_\-]+)+$`, val)

            // URL
            case "url":
                match = gregex.IsMatchString(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, val)

            // domain
            case "domain":
                match = gregex.IsMatchString(`^([0-9a-zA-Z][0-9a-zA-Z-]{0,62}\.)+([0-9a-zA-Z][0-9a-zA-Z-]{0,62})\.?$`, val)

            // IP(IPv4/IPv6)
            case "ip":
                match = gipv4.Validate(val) || gipv6.Validate(val)

            // IPv4
            case "ipv4":
                match = gipv4.Validate(val)

            // IPv6
            case "ipv6":
                match = gipv6.Validate(val)

            // MAC地址
            case "mac":
                match = gregex.IsMatchString(`^([0-9A-Fa-f]{2}[\-:]){5}[0-9A-Fa-f]{2}$`, val)

            default:
                errorMsgs[ruleKey] = "Invalid rule name:" + ruleKey
        }

        // 错误消息整合
        if !match {
            // 不存在则使用默认的错误信息，
            // 如果在校验过程中已经设置了错误信息，那么这里便不作处理
            if _, ok := errorMsgs[ruleKey]; !ok {
                if msg, ok := customMsgMap[ruleKey]; ok {
                    errorMsgs[ruleKey] = msg
                } else {
                    errorMsgs[ruleKey] = errorMsgMap.Get(ruleKey)
                }
            }
        }
        index++
    }
    if len(errorMsgs) > 0 {
        return newError([]string{rules}, ErrorMap {
            // 单条数值校验没有键名
            "" : errorMsgs,
        })
    }
    return nil
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
                if params[array[i]] != "" {
                    required = true
                    break
                }
            }

        // 必须字段(当所给定所有字段值都不为空时)
        case "required-with-all":
            required = true
            array   := strings.Split(ruleVal, ",")
            for i := 0; i < len(array); i++ {
                if params[array[i]] == "" {
                    required = false
                    break
                }
            }

        // 必须字段(当所给定任意字段值为空时)
        case "required-without":
            required = false
            array   := strings.Split(ruleVal, ",")
            for i := 0; i < len(array); i++ {
                if params[array[i]] == "" {
                    required = true
                    break
                }
            }

        // 必须字段(当所给定所有字段值都为空时)
        case "required-without-all":
            required = true
            array   := strings.Split(ruleVal, ",")
            for i := 0; i < len(array); i++ {
                if params[array[i]] != "" {
                    required = false
                    break
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
func checkLength(value, ruleKey, ruleVal string, customMsgMap map[string]string) string {
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
                if v, ok := customMsgMap[ruleKey]; !ok {
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
                    if v, ok := customMsgMap[ruleKey]; !ok {
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
                    if v, ok := customMsgMap[ruleKey]; !ok {
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
func checkSize(value, ruleKey, ruleVal string, customMsgMap map[string]string) string {
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
                    if v, ok := customMsgMap[ruleKey]; !ok {
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
                        if v, ok := customMsgMap[ruleKey]; !ok {
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
                        if v, ok := customMsgMap[ruleKey]; !ok {
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

