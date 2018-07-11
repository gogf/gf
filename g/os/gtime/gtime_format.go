// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtime

import (
    "strings"
)

var (
    // 参考：http://php.net/manual/zh/function.date.php
    formats = map[byte]string {
        // ================== 日 ==================
        'd'	 : "02",       // 月份中的第几天，有前导零的 2 位数字(01 到 31)
        'D'  : "Mon",      // 星期中的第几天，文本表示，3 个字母(Mon 到 Sun)
        'j'  : "=j=02",    // 月份中的第几天，没有前导零(1 到 31)
        'l'  : "Monday",   // ("L"的小写字母)星期几，完整的文本格式(Sunday 到 Saturday)

        // ================== 月 ==================
        'F'  : "January",  // 月份，完整的文本格式，例如 January 或者 March	January 到 December
        'm'  : "01",	   // 数字表示的月份，有前导零(01 到 12)
        'M'  : "Jan",      // 三个字母缩写表示的月份(Jan 到 Dec)
        'n'  : "1",        // 数字表示的月份，没有前导零(1 到 12)

        // ================== 年 ==================
        'Y'  : "2006",     // 4 位数字完整表示的年份, 例如：1999 或 2003
        'y'  : "06",       // 2 位数字表示的年份, 例如：99 或 03

        // ================== 时间 ==================
        'a'  : "pm",       // 小写的上午和下午值	am 或 pm
        'A'  : "PM",       // 大写的上午和下午值	AM 或 PM
        'g'  : "3",        // 小时，12 小时格式，没有前导零,  1 到 12
        'G'  : "=G=15",    // 小时，24 小时格式，没有前导零,  0 到 23
        'h'  : "03",       // 小时，12 小时格式，有前导零,   01 到 12
        'H'  : "15",       // 小时，24 小时格式，有前导零,   00 到 23
        'i'  : "04",       // 有前导零的分钟数, 00 到 59
        's'  : "05",       // 秒数，有前导零,   00 到 59
        'u'  : "=u=.000",  // 毫秒(3位)

        // ================== 时区 ==================
        'e'  : "MST",      // 时区标识, 例如：UTC，GMT，CST
        'O'  : "-0700",    // 与UTC相差的小时数, 例如：+0200
        'P'  : "-07:00",   // 与UTC的差别，小时和分钟之间有冒号分隔, 例如：+02:00
    }
)

// 将自定义的格式转换为标准库时间格式
func formatToStdLayout(format string) string {
    s := ""
    for i := 0; i < len(format); {
        switch format[i] {
            case '\\':
                if i < len(format)-1 {
                    s += string(format[i+1])
                    i += 2
                    continue
                } else {
                    return s
                }

            default:
                if f, ok := formats[format[i]]; ok {
                    switch format[i] {
                        case 'j':
                            s += "02"
                        case 'G':
                            s += "15"
                        case 'u':
                            if i > 0 && format[i - 1] == '.' {
                                s += "000"
                            } else {
                                s += ".000"
                            }
                        default:
                            s += f
                    }
                } else {
                    s += f
                }
        }
    }
    return s
}

// 格式化，使用自定义日期格式
func (t *Time) Format(format string) string {
    s := ""
    for i := 0; i < len(format); {
        switch format[i] {
            case '\\':
                if i < len(format) - 1 {
                    s += string(format[i + 1])
                    i += 2
                    continue
                } else {
                    return s
                }

            default:
                if f, ok := formats[format[i]]; ok {
                    r := t.Time.Format(f)
                    switch format[i] {
                        case 'j':
                            s += strings.Replace(r, "=j=0", "", -1)
                        case 'G':
                            s += strings.Replace(r, "=G=0", "", -1)
                        case 'u':
                            s += strings.Replace(r, "=u=.", "", -1)
                        default:
                            s += r
                    }
                } else {
                    s += f
                }
        }
        i++
    }
    return s
}

// 格式化，使用标准库格式
func (t *Time) Layout(layout string) string {
    return t.Time.Format(layout)
}