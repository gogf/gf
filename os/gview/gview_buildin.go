// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import (
	"fmt"
	"github.com/gogf/gf/util/gutil"
	"strings"

	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"

	htmltpl "html/template"
)

// funcDump implements build-in template function: dump
func (view *View) funcDump(values ...interface{}) (result string) {
	result += "<!--\n"
	for _, v := range values {
		result += gutil.Export(v) + "\n"
	}
	result += "-->\n"
	return result
}

// funcEq implements build-in template function: eq
func (view *View) funcEq(value interface{}, others ...interface{}) bool {
	s := gconv.String(value)
	for _, v := range others {
		if strings.Compare(s, gconv.String(v)) != 0 {
			return false
		}
	}
	return true
}

// funcNe implements build-in template function: ne
func (view *View) funcNe(value, other interface{}) bool {
	return strings.Compare(gconv.String(value), gconv.String(other)) != 0
}

// funcLt implements build-in template function: lt
func (view *View) funcLt(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) < gconv.Int64(other)
	}
	return strings.Compare(s1, s2) < 0
}

// funcLe implements build-in template function: le
func (view *View) funcLe(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) <= gconv.Int64(other)
	}
	return strings.Compare(s1, s2) <= 0
}

// funcGt implements build-in template function: gt
func (view *View) funcGt(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) > gconv.Int64(other)
	}
	return strings.Compare(s1, s2) > 0
}

// funcGe implements build-in template function: ge
func (view *View) funcGe(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) >= gconv.Int64(other)
	}
	return strings.Compare(s1, s2) >= 0
}

// funcInclude implements build-in template function: include
// Note that configuration AutoEncode does not affect the output of this function.
func (view *View) funcInclude(file interface{}, data ...map[string]interface{}) htmltpl.HTML {
	var m map[string]interface{} = nil
	if len(data) > 0 {
		m = data[0]
	}
	path := gconv.String(file)
	if path == "" {
		return ""
	}
	// It will search the file internally.
	content, err := view.Parse(path, m)
	if err != nil {
		return htmltpl.HTML(err.Error())
	}
	return htmltpl.HTML(content)
}

// funcText implements build-in template function: text
func (view *View) funcText(html interface{}) string {
	return ghtml.StripTags(gconv.String(html))
}

// funcHtmlEncode implements build-in template function: html
func (view *View) funcHtmlEncode(html interface{}) string {
	return ghtml.Entities(gconv.String(html))
}

// funcHtmlDecode implements build-in template function: htmldecode
func (view *View) funcHtmlDecode(html interface{}) string {
	return ghtml.EntitiesDecode(gconv.String(html))
}

// funcUrlEncode implements build-in template function: url
func (view *View) funcUrlEncode(url interface{}) string {
	return gurl.Encode(gconv.String(url))
}

// funcUrlDecode implements build-in template function: urldecode
func (view *View) funcUrlDecode(url interface{}) string {
	if content, err := gurl.Decode(gconv.String(url)); err == nil {
		return content
	} else {
		return err.Error()
	}
}

// funcDate implements build-in template function: date
func (view *View) funcDate(format interface{}, timestamp ...interface{}) string {
	t := int64(0)
	if len(timestamp) > 0 {
		t = gconv.Int64(timestamp[0])
	}
	if t == 0 {
		t = gtime.Timestamp()
	}
	return gtime.NewFromTimeStamp(t).Format(gconv.String(format))
}

// funcCompare implements build-in template function: compare
func (view *View) funcCompare(value1, value2 interface{}) int {
	return strings.Compare(gconv.String(value1), gconv.String(value2))
}

// funcSubStr implements build-in template function: substr
func (view *View) funcSubStr(start, end, str interface{}) string {
	return gstr.SubStrRune(gconv.String(str), gconv.Int(start), gconv.Int(end))
}

// funcStrLimit implements build-in template function: strlimit
func (view *View) funcStrLimit(length, suffix, str interface{}) string {
	return gstr.StrLimitRune(gconv.String(str), gconv.Int(length), gconv.String(suffix))
}

// funcConcat implements build-in template function: concat
func (view *View) funcConcat(str ...interface{}) string {
	var s string
	for _, v := range str {
		s += gconv.String(v)
	}
	return s
}

// funcReplace implements build-in template function: replace
func (view *View) funcReplace(search, replace, str interface{}) string {
	return gstr.Replace(gconv.String(str), gconv.String(search), gconv.String(replace), -1)
}

// funcHighlight implements build-in template function: highlight
func (view *View) funcHighlight(key, color, str interface{}) string {
	return gstr.Replace(gconv.String(str), gconv.String(key), fmt.Sprintf(`<span style="color:%v;">%v</span>`, color, key))
}

// funcHideStr implements build-in template function: hidestr
func (view *View) funcHideStr(percent, hide, str interface{}) string {
	return gstr.HideStr(gconv.String(str), gconv.Int(percent), gconv.String(hide))
}

// funcToUpper implements build-in template function: toupper
func (view *View) funcToUpper(str interface{}) string {
	return gstr.ToUpper(gconv.String(str))
}

// funcToLower implements build-in template function: toupper
func (view *View) funcToLower(str interface{}) string {
	return gstr.ToLower(gconv.String(str))
}

// funcNl2Br implements build-in template function: nl2br
func (view *View) funcNl2Br(str interface{}) string {
	return gstr.Nl2Br(gconv.String(str))
}
