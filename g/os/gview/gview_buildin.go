// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import (
	"fmt"
	"github.com/gogf/gf/g/encoding/ghtml"
	"github.com/gogf/gf/g/encoding/gurl"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
	"strings"
)

// Build-in template function: eq
func (view *View) funcEq(value interface{}, others...interface{}) bool {
	s := gconv.String(value)
	for _, v := range others {
		if strings.Compare(s, gconv.String(v)) != 0 {
			return false
		}
	}
	return true
}

// Build-in template function: ne
func (view *View) funcNe(value interface{}, other interface{}) bool {
	return strings.Compare(gconv.String(value), gconv.String(other)) != 0
}

// Build-in template function: lt
func (view *View) funcLt(value interface{}, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) < gconv.Int64(other)
	}
	return strings.Compare(s1, s2) < 0
}

// Build-in template function: le
func (view *View) funcLe(value interface{}, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) <= gconv.Int64(other)
	}
	return strings.Compare(s1, s2) <= 0
}

// Build-in template function: gt
func (view *View) funcGt(value interface{}, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) > gconv.Int64(other)
	}
	return strings.Compare(s1, s2) > 0
}

// Build-in template function: ge
func (view *View) funcGe(value interface{}, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) >= gconv.Int64(other)
	}
	return strings.Compare(s1, s2) >= 0
}

// Build-in template function: include
func (view *View) funcInclude(file string, data...map[string]interface{}) string {
    var m map[string]interface{} = nil
    if len(data) > 0 {
        m = data[0]
    }
    content, err := view.Parse(file, m)
    if err != nil {
        return err.Error()
    }
    return content
}

// Build-in template function: text
func (view *View) funcText(html interface{}) string {
    return ghtml.StripTags(gconv.String(html))
}

// Build-in template function: html
func (view *View) funcHtmlEncode(html interface{}) string {
    return ghtml.Entities(gconv.String(html))
}

// Build-in template function: htmldecode
func (view *View) funcHtmlDecode(html interface{}) string {
    return ghtml.EntitiesDecode(gconv.String(html))
}

// Build-in template function: url
func (view *View) funcUrlEncode(url interface{}) string {
    return gurl.Encode(gconv.String(url))
}

// Build-in template function: urldecode
func (view *View) funcUrlDecode(url interface{}) string {
    if content, err := gurl.Decode(gconv.String(url)); err == nil {
        return content
    } else {
        return err.Error()
    }
}

// Build-in template function: date
func (view *View) funcDate(format string, timestamp...interface{}) string {
    t := int64(0)
    if len(timestamp) > 0 {
        t = gconv.Int64(timestamp[0])
    }
    if t == 0 {
        t = gtime.Millisecond()
    }
    return gtime.NewFromTimeStamp(t).Format(format)
}

// Build-in template function: compare
func (view *View) funcCompare(value1, value2 interface{}) int {
    return strings.Compare(gconv.String(value1), gconv.String(value2))
}

// Build-in template function: substr
func (view *View) funcSubStr(start, end int, str interface{}) string {
    return gstr.SubStr(gconv.String(str), start, end)
}

// Build-in template function: strlimit
func (view *View) funcStrLimit(length int, suffix string, str interface{}) string {
    return gstr.StrLimit(gconv.String(str), length, suffix)
}

// Build-in template function: highlight
func (view *View) funcHighlight(key string, color string, str interface{}) string {
    return gstr.Replace(gconv.String(str), key, fmt.Sprintf(`<span style="color:%s;">%s</span>`, color, key))
}

// Build-in template function: hidestr
func (view *View) funcHideStr(percent int, hide string, str interface{}) string {
    return gstr.HideStr(gconv.String(str), percent, hide)
}

// Build-in template function: toupper
func (view *View) funcToUpper(str interface{}) string {
    return gstr.ToUpper(gconv.String(str))
}

// Build-in template function: toupper
func (view *View) funcToLower(str interface{}) string {
    return gstr.ToLower(gconv.String(str))
}

// Build-in template function: nl2br
func (view *View) funcNl2Br(str interface{}) string {
    return gstr.Nl2Br(gconv.String(str))
}


