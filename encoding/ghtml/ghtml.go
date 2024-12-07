// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package ghtml provides useful API for HTML content handling.
package ghtml

import (
	"html"
	"reflect"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// StripTags strips HTML tags from content, and returns only text.
// Referer: http://php.net/manual/zh/function.strip-tags.php
func StripTags(s string) string {
	return strip.StripTags(s)
}

// Entities encodes all HTML chars for content.
// Referer: http://php.net/manual/zh/function.htmlentities.php
func Entities(s string) string {
	return html.EscapeString(s)
}

// EntitiesDecode decodes all HTML chars for content.
// Referer: http://php.net/manual/zh/function.html-entity-decode.php
func EntitiesDecode(s string) string {
	return html.UnescapeString(s)
}

// SpecialChars encodes some special chars for content, these special chars are:
// "&", "<", ">", `"`, "'".
// Referer: http://php.net/manual/zh/function.htmlspecialchars.php
func SpecialChars(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&#34;",
		"'", "&#39;",
	).Replace(s)
}

// SpecialCharsDecode decodes some special chars for content, these special chars are:
// "&", "<", ">", `"`, "'".
// Referer: http://php.net/manual/zh/function.htmlspecialchars-decode.php
func SpecialCharsDecode(s string) string {
	return strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&#34;", `"`,
		"&#39;", "'",
	).Replace(s)
}

// SpecialCharsMapOrStruct automatically encodes string values/attributes for map/struct.
//
// Note that, if operation on struct, the given parameter `mapOrStruct` should be type of pointer to struct.
//
// For example:
// var m = map{}
// var s = struct{}{}
// OK: SpecialCharsMapOrStruct(m)
// OK: SpecialCharsMapOrStruct(&s)
// Error: SpecialCharsMapOrStruct(s)
func SpecialCharsMapOrStruct(mapOrStruct interface{}) error {
	var (
		reflectValue = reflect.ValueOf(mapOrStruct)
		reflectKind  = reflectValue.Kind()
		originalKind = reflectKind
	)
	for reflectValue.IsValid() && (reflectKind == reflect.Ptr || reflectKind == reflect.Interface) {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}

	switch reflectKind {
	case reflect.Map:
		var (
			mapKeys  = reflectValue.MapKeys()
			mapValue reflect.Value
		)
		for _, key := range mapKeys {
			mapValue = reflectValue.MapIndex(key)
			switch mapValue.Kind() {
			case reflect.String:
				reflectValue.SetMapIndex(key, reflect.ValueOf(SpecialChars(mapValue.String())))
			case reflect.Interface:
				if mapValue.Elem().Kind() == reflect.String {
					reflectValue.SetMapIndex(
						key,
						reflect.ValueOf(SpecialChars(mapValue.Elem().String())),
					)
				}
			default:
			}
		}

	case reflect.Struct:
		if originalKind != reflect.Ptr {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid input parameter type "%s", should be type of pointer to struct`,
				reflect.TypeOf(mapOrStruct).String(),
			)
		}
		var fieldValue reflect.Value
		for i := 0; i < reflectValue.NumField(); i++ {
			fieldValue = reflectValue.Field(i)
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.Set(
					reflect.ValueOf(
						SpecialChars(fieldValue.String()),
					),
				)
			default:
			}
		}

	default:
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid input parameter type "%s"`,
			reflect.TypeOf(mapOrStruct).String(),
		)
	}
	return nil
}
