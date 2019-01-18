// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gtime provides functionality for measuring and displaying time.
// 
// 时间管理.
package gtime

import (
    "errors"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/util/gstr"
    "regexp"
    "strconv"
    "strings"
    "time"
)

const (
    // 时间间隔缩写
    D  = 24*time.Hour
    H  = time.Hour
    M  = time.Minute
    S  = time.Second
    MS = time.Millisecond
    US = time.Microsecond
    NS = time.Nanosecond

    // 常用时间格式正则匹配，支持的标准时间格式：
    // "2017-12-14 04:51:34 +0805 LMT",
    // "2017-12-14 04:51:34 +0805 LMT",
    // "2006-01-02T15:04:05Z07:00",
    // "2014-01-17T01:19:15+08:00",
    // "2018-02-09T20:46:17.897Z",
    // "2018-02-09 20:46:17.897",
    // "2018-02-09T20:46:17Z",
    // "2018-02-09 20:46:17",
    // "2018/10/31 - 16:38:46"
    // "2018-02-09",
    // "2018.02.09",
    // 日期连接符号支持'-'、'/'、'.'
    TIME_REAGEX_PATTERN1 = `(\d{4}[-/\.]\d{2}[-/\.]\d{2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`
    // 01-Nov-2018 11:50:28
    // 01/Nov/2018 11:50:28
    // 01.Nov.2018 11:50:28
    // 01.Nov.2018:11:50:28
    // 日期连接符号支持'-'、'/'、'.'
    TIME_REAGEX_PATTERN2 = `(\d{1,2}[-/\.][A-Za-z]{3,}[-/\.]\d{4})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`
)

var (
    // 使用正则判断会比直接使用ParseInLocation挨个轮训判断要快很多
    timeRegex1, _   = regexp.Compile(TIME_REAGEX_PATTERN1)
    timeRegex2, _   = regexp.Compile(TIME_REAGEX_PATTERN2)
    // 月份英文与阿拉伯数字对应关系
    monthMap = map[string]int {
        "jan"        : 1,
        "feb"        : 2,
        "mar"        : 3,
        "apr"        : 4,
        "may"        : 5,
        "jun"        : 6,
        "jul"        : 7,
        "aug"        : 8,
        "sep"        : 9,
        "sept"       : 9,
        "oct"        : 10,
        "nov"        : 11,
        "dec"        : 12,
        "january"    : 1,
        "february"   : 2,
        "march"      : 3,
        "april"      : 4,
        "june"       : 6,
        "july"       : 7,
        "august"     : 8,
        "september"  : 9,
        "october"    : 10,
        "november"   : 11,
        "december"   : 12,
    }
)

// 设置当前进程全局的默认时区，如: Asia/Shanghai
func SetTimeZone(zone string) error {
    location, err := time.LoadLocation(zone)
    if err == nil {
        time.Local = location
    }
    return err
}

// 获取当前的纳秒数
func Nanosecond() int64 {
    return time.Now().UnixNano()
}

// 获取当前的微秒数
func Microsecond() int64 {
    return time.Now().UnixNano()/1e3
}

// 获取当前的毫秒数
func Millisecond() int64 {
    return time.Now().UnixNano()/1e6
}

// 获取当前的秒数(时间戳)
func Second() int64 {
    return time.Now().Unix()
}

// 获得当前的日期(例如：2006-01-02)
func Date() string {
    return time.Now().Format("2006-01-02")
}

// 获得当前的时间(例如：2006-01-02 15:04:05)
func Datetime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

// 解析日期字符串(日期支持'-'或'/'或'.'连接符号)
func parseDateStr(s string) (year, month, day int) {
    array := strings.Split(s, "-")
    if len(array) < 3 {
        array = strings.Split(s, "/")
    }
    if len(array) < 3 {
        array = strings.Split(s, ".")
    }
    // 解析失败
    if len(array) < 3 {
        return
    }
    // 判断年份在开头还是末尾
    if gstr.IsNumeric(array[1]) {
        year, _  = strconv.Atoi(array[0])
        month, _ = strconv.Atoi(array[1])
        day, _   = strconv.Atoi(array[2])
    } else {
        if v, ok := monthMap[strings.ToLower(array[1])]; ok {
            month = v
        } else {
            return
        }
        year, _  = strconv.Atoi(array[2])
        day, _   = strconv.Atoi(array[1])
    }
    // 年是否为缩写，如果是，那么需要补上前缀
    if year < 100 {
        year = int(time.Now().Year()/100)*100 + year
    }
    return
}

// 字符串转换为时间对象，format参数指定格式的format(如: Y-m-d H:i:s)，当指定format参数时效果同StrToTimeFormat方法。
// 注意：自动解析日期时间时，必须有日期才能解析成功，如果字符串中不带有日期字段，那么解析失败。
func StrToTime(str string, format...string) (*Time, error) {
    if len(format) > 0 {
        return StrToTimeFormat(str, format[0])
    }
    var year, month, day int
    var hour, min, sec, nsec int
    var match []string
    var local = time.Local
    if match = timeRegex1.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
        for k, v := range match {
            match[k] = strings.TrimSpace(v)
        }
        year, month, day = parseDateStr(match[1])
    } else if match = timeRegex2.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
        for k, v := range match {
            match[k] = strings.TrimSpace(v)
        }
        year, month, day = parseDateStr(match[1])
    } else {
        return nil, errors.New("unsupported time format")
    }

    // 时间
    if len(match[2]) > 0 {
        s := strings.Replace(match[2], ":", "", -1)
        if len(s) < 6 {
            s += strings.Repeat("0", 6 - len(s))
        }
        hour, _ = strconv.Atoi(s[0 : 2])
        min, _  = strconv.Atoi(s[2 : 4])
        sec, _  = strconv.Atoi(s[4 : 6])
    }
    // 纳秒，检查并执行位补齐
    if len(match[3]) > 0 {
        nsec, _   = strconv.Atoi(match[3])
        for i := 0; i < 9 - len(match[3]); i++ {
            nsec *= 10
        }
    }
    // 如果字符串中有时区信息(具体时间信息)，那么执行时区转换，将时区转成UTC
    if match[4] != "" && match[6] == "" {
        match[6] = "000000"
    }
    // 如果offset有值优先处理offset，否则处理后面的时区名称
    if match[6] != "" {
        zone := strings.Replace(match[6], ":", "", -1)
        zone  = strings.TrimLeft(zone, "+-")
        if len(zone) <= 6 {
            zone += strings.Repeat("0", 6 - len(zone))
            h, _ := strconv.Atoi(zone[0 : 2])
            m, _ := strconv.Atoi(zone[2 : 4])
            s, _ := strconv.Atoi(zone[4 : 6])
            // 判断字符串输入的时区是否和当前程序时区相等(使用offset判断)，不相等则将对象统一转换为UTC时区
            // 当前程序时区Offset(秒)
            _, localOffset := time.Now().Zone()
            if (h * 3600 + m * 60 + s) != localOffset {
                local = time.UTC
                // UTC时差转换
                operation := match[5]
                if operation != "+" && operation != "-" {
                    operation = "-"
                }
                switch operation {
                case "+":
                    if h > 0 {
                        hour -= h
                    }
                    if m > 0 {
                        min  -= m
                    }
                    if s > 0 {
                        sec  -= s
                    }
                case "-":
                    if h > 0 {
                        hour += h
                    }
                    if m > 0 {
                        min  += m
                    }
                    if s > 0 {
                        sec  += s
                    }
                }
            }
        }
    }
    // 统一生成UTC时间对象
    return NewFromTime(time.Date(year, time.Month(month), day, hour, min, sec, nsec, local)), nil
}

// 时区转换
func ConvertZone(strTime string, toZone string, fromZone...string) (*Time, error) {
   t, err := StrToTime(strTime)
   if err != nil {
       return nil, err
   }
   if len(fromZone) > 0 {
       if l, err := time.LoadLocation(fromZone[0]); err != nil {
           return nil, err
       } else {
           t.Time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Time.Second(), t.Time.Nanosecond(), l)
       }
   }
    if l, err := time.LoadLocation(toZone); err != nil {
        return nil, err
    } else {
        return t.ToLocation(l), nil
    }
}

// 字符串转换为时间对象，指定字符串时间格式，format格式形如：Y-m-d H:i:s
func StrToTimeFormat(str string, format string) (*Time, error) {
    return StrToTimeLayout(str, formatToStdLayout(format))
}

// 字符串转换为时间对象，通过标准库layout格式进行解析，layout格式形如：2006-01-02 15:04:05
func StrToTimeLayout(str string, layout string) (*Time, error) {
    if t, err := time.ParseInLocation(layout, str, time.Local); err == nil {
        return NewFromTime(t), nil
    } else {
        return nil, err
    }
}

// 从字符串内容中(也可以是文件名称等等)解析时间，并返回解析成功的时间对象，否则返回nil。
// 注意当内容中存在多个时间时，会解析第一个。
// format参数可以指定需要解析的时间格式。
func ParseTimeFromContent(content string, format...string) *Time {
    if len(format) > 0 {
        if match, err := gregex.MatchString(formatToRegexPattern(format[0]), content); err == nil && len(match) > 0 {
            return NewFromStrFormat(match[0], format[0])
        }
    } else {
        if match := timeRegex1.FindStringSubmatch(content); len(match) >= 1 {
            return NewFromStr(strings.Trim(match[0], "./_- \n\r"))
        } else if match := timeRegex2.FindStringSubmatch(content); len(match) >= 1 {
            return NewFromStr(strings.Trim(match[0], "./_- \n\r"))
        }
    }
    return nil
}

// 计算函数f执行的时间，单位纳秒
func FuncCost(f func()) int64 {
    t := Nanosecond()
    f()
    return Nanosecond() - t
}