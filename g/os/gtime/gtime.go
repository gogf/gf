// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 时间管理
package gtime

import (
    "time"
    "regexp"
    "strings"
    "strconv"
    "errors"
)

const (
    TIME_REAGEX_PATTERN = `(\d{4}[-/]\d{2}[-/]\d{2})[\sT]{0,1}(\d{2}:\d{2}:\d{2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`
)

var (
    // 用于time.Time转换使用，防止多次Compile
    timeRegex   *regexp.Regexp
)

func init() {
    // 使用正则判断会比直接使用ParseInLocation挨个轮训判断要快很多
    timeRegex, _   = regexp.Compile(TIME_REAGEX_PATTERN)

}

// 类似与js中的SetTimeout，一段时间后执行回调函数
func SetTimeout(t time.Duration, callback func()) {
    go func() {
        time.Sleep(t)
        callback()
    }()
}

// 类似与js中的SetInterval，每隔一段时间后执行回调函数，当回调函数返回true，那么继续执行，否则终止执行，该方法是异步的
// 注意：由于采用的是循环而不是递归操作，因此间隔时间将会以上一次回调函数执行完成的时间来计算
func SetInterval(t time.Duration, callback func() bool) {
    go func() {
        for {
            time.Sleep(t)
            if !callback() {
                break
            }
        }
    }()
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
    return time.Now().UnixNano()/1e9
}

// 获得当前的日期(例如：2006-01-02)
func Date() string {
    return time.Now().Format("2006-01-02")
}

// 获得当前的时间(例如：2006-01-02 15:04:05)
func Datetime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

// 时间戳转换为指定格式的字符串，format格式形如：2006-01-02 03:04:05 PM
// 第二个参数指定需要格式化的时间戳，为非必需参数，默认为当前时间戳
func Format(format string, timestamps...int64) string {
    timestamp := Second()
    if len(timestamps) > 0 {
        timestamp = timestamps[0]
    }
    return time.Unix(timestamp, 0).Format(format)
}

// 字符串转换为时间对象，需要给定字符串时间格式，format格式形如：2006-01-02 15:04:05
// 不传递自定义格式下默认支持的标准时间格式：
// "2017-12-14 04:51:34 +0805 LMT",
// "2006-01-02T15:04:05Z07:00",
// "2014-01-17T01:19:15+08:00",
// "2018-02-09T20:46:17.897Z",
// "2018-02-09 20:46:17.897",
// "2018-02-09T20:46:17Z",
// "2018-02-09 20:46:17",
// "2018-02-09",
func StrToTime(str string, format...string) (time.Time, error) {
    // 优先使用用户输入日期格式进行转换
    if len(format) > 0 {
        if t, err := time.ParseInLocation(format[0], str, time.Local); err == nil {
            return t, nil
        } else {
            return time.Time{}, err
        }
    }
    var result time.Time
    var local  = time.Local
    if match := timeRegex.FindStringSubmatch(str); len(match) > 0 {
        var year, month, day, hour, min, sec, nsec int
        var array []string
        // 日期
        array = strings.Split(match[1], "-")
        if len(array) >= 3 {
            year, _  = strconv.Atoi(array[0])
            month, _ = strconv.Atoi(array[1])
            day, _   = strconv.Atoi(array[2])
        }
        // 时间
        array = strings.Split(match[2], ":")
        if len(array) >= 3 {
            hour, _  = strconv.Atoi(array[0])
            min, _   = strconv.Atoi(array[1])
            sec, _   = strconv.Atoi(array[2])
        }
        array = strings.Split(match[1], "-")
        // 纳秒，检查病执行位补齐
        if match[3] != "" {
            nsec, _   = strconv.Atoi(match[3])
            for i := 0; i < 9 - len(match[3]); i++ {
                nsec *= 10
            }
        }
        // 如果字符串中有时区信息，那么执行时区转换，将时区转成UTC
        if match[4] != "" && match[6] == "" {
            match[6] = "000000"
        }
        if match[6] != "" {
            zone := strings.Replace(match[6], ":", "", -1)
            zone  = strings.TrimLeft(zone, "+-")
            zone += strings.Repeat("0", 6 - len(zone))
            h, _ := strconv.Atoi(zone[0 : 2])
            m, _ := strconv.Atoi(zone[2 : 4])
            s, _ := strconv.Atoi(zone[4 : 6])
            // 判断字符串输入的时区是否和当前程序时区相等，不相等则将对象统一转换为UTC时区
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
        // 生成UTC时间对象
        result = time.Date(year, time.Month(month), day, hour, min, sec, nsec, local)
        return result, nil
    }
    return result, errors.New("unsupported time format")
}
