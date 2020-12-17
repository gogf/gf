// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import (
	"fmt"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/util/gutil"
	"strings"

	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"

	htmltpl "html/template"
)

// buildInFuncDump implements build-in template function: dump
func (view *View) buildInFuncDump(values ...interface{}) (result string) {
	result += "<!--\n"
	for _, v := range values {
		result += gutil.Export(v) + "\n"
	}
	result += "-->\n"
	return result
}

// buildInFuncMap implements build-in template function: map
func (view *View) buildInFuncMap(value ...interface{}) map[string]interface{} {
	if len(value) > 0 {
		return gconv.Map(value[0])
	}
	return map[string]interface{}{}
}

// buildInFuncMaps implements build-in template function: maps
func (view *View) buildInFuncMaps(value ...interface{}) []map[string]interface{} {
	if len(value) > 0 {
		return gconv.Maps(value[0])
	}
	return []map[string]interface{}{}
}

// buildInFuncEq implements build-in template function: eq
func (view *View) buildInFuncEq(value interface{}, others ...interface{}) bool {
	s := gconv.String(value)
	for _, v := range others {
		if strings.Compare(s, gconv.String(v)) != 0 {
			return false
		}
	}
	return true
}

// buildInFuncNe implements build-in template function: ne
func (view *View) buildInFuncNe(value, other interface{}) bool {
	return strings.Compare(gconv.String(value), gconv.String(other)) != 0
}

// buildInFuncLt implements build-in template function: lt
func (view *View) buildInFuncLt(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) < gconv.Int64(other)
	}
	return strings.Compare(s1, s2) < 0
}

// buildInFuncLe implements build-in template function: le
func (view *View) buildInFuncLe(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) <= gconv.Int64(other)
	}
	return strings.Compare(s1, s2) <= 0
}

// buildInFuncGt implements build-in template function: gt
func (view *View) buildInFuncGt(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) > gconv.Int64(other)
	}
	return strings.Compare(s1, s2) > 0
}

// buildInFuncGe implements build-in template function: ge
func (view *View) buildInFuncGe(value, other interface{}) bool {
	s1 := gconv.String(value)
	s2 := gconv.String(other)
	if gstr.IsNumeric(s1) && gstr.IsNumeric(s2) {
		return gconv.Int64(value) >= gconv.Int64(other)
	}
	return strings.Compare(s1, s2) >= 0
}

// buildInFuncInclude implements build-in template function: include
// Note that configuration AutoEncode does not affect the output of this function.
func (view *View) buildInFuncInclude(file interface{}, data ...map[string]interface{}) htmltpl.HTML {
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

// buildInFuncText implements build-in template function: text
func (view *View) buildInFuncText(html interface{}) string {
	return ghtml.StripTags(gconv.String(html))
}

// buildInFuncHtmlEncode implements build-in template function: html
func (view *View) buildInFuncHtmlEncode(html interface{}) string {
	return ghtml.Entities(gconv.String(html))
}

// buildInFuncHtmlDecode implements build-in template function: htmldecode
func (view *View) buildInFuncHtmlDecode(html interface{}) string {
	return ghtml.EntitiesDecode(gconv.String(html))
}

// buildInFuncUrlEncode implements build-in template function: url
func (view *View) buildInFuncUrlEncode(url interface{}) string {
	return gurl.Encode(gconv.String(url))
}

// buildInFuncUrlDecode implements build-in template function: urldecode
func (view *View) buildInFuncUrlDecode(url interface{}) string {
	if content, err := gurl.Decode(gconv.String(url)); err == nil {
		return content
	} else {
		return err.Error()
	}
}

// buildInFuncDate implements build-in template function: date
func (view *View) buildInFuncDate(format interface{}, timestamp ...interface{}) string {
	t := int64(0)
	if len(timestamp) > 0 {
		t = gconv.Int64(timestamp[0])
	}
	if t == 0 {
		t = gtime.Timestamp()
	}
	return gtime.NewFromTimeStamp(t).Format(gconv.String(format))
}

// buildInFuncCompare implements build-in template function: compare
func (view *View) buildInFuncCompare(value1, value2 interface{}) int {
	return strings.Compare(gconv.String(value1), gconv.String(value2))
}

// buildInFuncSubStr implements build-in template function: substr
func (view *View) buildInFuncSubStr(start, end, str interface{}) string {
	return gstr.SubStrRune(gconv.String(str), gconv.Int(start), gconv.Int(end))
}

// buildInFuncStrLimit implements build-in template function: strlimit
func (view *View) buildInFuncStrLimit(length, suffix, str interface{}) string {
	return gstr.StrLimitRune(gconv.String(str), gconv.Int(length), gconv.String(suffix))
}

// buildInFuncConcat implements build-in template function: concat
func (view *View) buildInFuncConcat(str ...interface{}) string {
	var s string
	for _, v := range str {
		s += gconv.String(v)
	}
	return s
}

// buildInFuncReplace implements build-in template function: replace
func (view *View) buildInFuncReplace(search, replace, str interface{}) string {
	return gstr.Replace(gconv.String(str), gconv.String(search), gconv.String(replace), -1)
}

// buildInFuncHighlight implements build-in template function: highlight
func (view *View) buildInFuncHighlight(key, color, str interface{}) string {
	return gstr.Replace(gconv.String(str), gconv.String(key), fmt.Sprintf(`<span style="color:%v;">%v</span>`, color, key))
}

// buildInFuncHideStr implements build-in template function: hidestr
func (view *View) buildInFuncHideStr(percent, hide, str interface{}) string {
	return gstr.HideStr(gconv.String(str), gconv.Int(percent), gconv.String(hide))
}

// buildInFuncToUpper implements build-in template function: toupper
func (view *View) buildInFuncToUpper(str interface{}) string {
	return gstr.ToUpper(gconv.String(str))
}

// buildInFuncToLower implements build-in template function: toupper
func (view *View) buildInFuncToLower(str interface{}) string {
	return gstr.ToLower(gconv.String(str))
}

// buildInFuncNl2Br implements build-in template function: nl2br
func (view *View) buildInFuncNl2Br(str interface{}) string {
	return gstr.Nl2Br(gconv.String(str))
}

// buildInFuncJson implements build-in template function: json ,
// which encodes and returns <value> as JSON string.
func (view *View) buildInFuncJson(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	return gconv.UnsafeBytesToStr(b), err
}
