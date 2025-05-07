// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import (
	"bytes"
	"context"
	"fmt"
	htmltpl "html/template"
	"strings"

	"github.com/gogf/gf/v3/encoding/ghtml"
	"github.com/gogf/gf/v3/encoding/gjson"
	"github.com/gogf/gf/v3/encoding/gurl"
	"github.com/gogf/gf/v3/os/gtime"
	"github.com/gogf/gf/v3/text/gstr"
	"github.com/gogf/gf/v3/util/gconv"
	"github.com/gogf/gf/v3/util/gmode"
	"github.com/gogf/gf/v3/util/gutil"
)

// buildInFuncDump implements build-in template function: dump
func (view *View) buildInFuncDump(values ...interface{}) string {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("\n")
	buffer.WriteString("<!--\n")
	if gmode.IsDevelop() {
		for _, v := range values {
			gutil.DumpTo(buffer, v, gutil.DumpOption{})
			buffer.WriteString("\n")
		}
	} else {
		buffer.WriteString("dump feature is disabled as process is not running in develop mode\n")
	}
	buffer.WriteString("-->\n")
	return buffer.String()
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
		if strings.Compare(s, gconv.String(v)) == 0 {
			return true
		}
	}
	return false
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
	content, err := view.Parse(context.TODO(), path, m)
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
	return gtime.NewFromTimeStamp(t).Layout(gconv.String(format))
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
// which encodes and returns `value` as JSON string.
func (view *View) buildInFuncJson(value interface{}) (string, error) {
	b, err := gjson.Marshal(value)
	return string(b), err
}

// buildInFuncXml implements build-in template function: xml ,
// which encodes and returns `value` as XML string.
func (view *View) buildInFuncXml(value interface{}, rootTag ...string) (string, error) {
	b, err := gjson.New(value).ToXml(rootTag...)
	return string(b), err
}

// buildInFuncIni implements build-in template function: ini ,
// which encodes and returns `value` as XML string.
func (view *View) buildInFuncIni(value interface{}) (string, error) {
	b, err := gjson.New(value).ToIni()
	return string(b), err
}

// buildInFuncYaml implements build-in template function: yaml ,
// which encodes and returns `value` as YAML string.
func (view *View) buildInFuncYaml(value interface{}) (string, error) {
	b, err := gjson.New(value).ToYaml()
	return string(b), err
}

// buildInFuncYamlIndent implements build-in template function: yamli ,
// which encodes and returns `value` as YAML string with custom indent string.
func (view *View) buildInFuncYamlIndent(value, indent interface{}) (string, error) {
	b, err := gjson.New(value).ToYamlIndent(gconv.String(indent))
	return string(b), err
}

// buildInFuncToml implements build-in template function: toml ,
// which encodes and returns `value` as TOML string.
func (view *View) buildInFuncToml(value interface{}) (string, error) {
	b, err := gjson.New(value).ToToml()
	return string(b), err
}

// buildInFuncPlus implements build-in template function: plus ,
// which returns the result that pluses all `deltas` to `value`.
func (view *View) buildInFuncPlus(value interface{}, deltas ...interface{}) string {
	result := gconv.Float64(value)
	for _, v := range deltas {
		result += gconv.Float64(v)
	}
	return gconv.String(result)
}

// buildInFuncMinus implements build-in template function: minus ,
// which returns the result that subtracts all `deltas` from `value`.
func (view *View) buildInFuncMinus(value interface{}, deltas ...interface{}) string {
	result := gconv.Float64(value)
	for _, v := range deltas {
		result -= gconv.Float64(v)
	}
	return gconv.String(result)
}

// buildInFuncTimes implements build-in template function: times ,
// which returns the result that multiplies `value` by all of `values`.
func (view *View) buildInFuncTimes(value interface{}, values ...interface{}) string {
	result := gconv.Float64(value)
	for _, v := range values {
		result *= gconv.Float64(v)
	}
	return gconv.String(result)
}

// buildInFuncDivide implements build-in template function: divide ,
// which returns the result that divides `value` by all of `values`.
func (view *View) buildInFuncDivide(value interface{}, values ...interface{}) string {
	result := gconv.Float64(value)
	for _, v := range values {
		value2Float64 := gconv.Float64(v)
		if value2Float64 == 0 {
			// Invalid `value2`.
			return "0"
		}
		result /= value2Float64
	}
	return gconv.String(result)
}
