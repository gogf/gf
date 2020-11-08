// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package ghtml provides useful API for HTML content handling.
package ghtml

import (
	"html"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
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
