// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package ghtml provides useful API for HTML content handling.
//
// HTML编码.
package ghtml

import (
    "strings"
    "html"
    "gitee.com/johng/gf/third/github.com/grokify/html-strip-tags-go"
)

// 过滤掉HTML标签，只返回text内容
// 参考：http://php.net/manual/zh/function.strip-tags.php
func StripTags(s string) string {
    return strip.StripTags(s)
}

// 本函数各方面都和SpecialChars一样，
// 除了Entities会转换所有具有 HTML 实体的字符。
// 参考：http://php.net/manual/zh/function.htmlentities.php
func Entities(s string) string {
    return html.EscapeString(s)
}

// Entities 的相反操作
// 参考：http://php.net/manual/zh/function.html-entity-decode.php
func EntitiesDecode(s string) string {
    return html.UnescapeString(s)
}

// 将html中的部分特殊标签转换为html转义标签
// 参考：http://php.net/manual/zh/function.htmlspecialchars.php
func SpecialChars(s string) string {
    return strings.NewReplacer(
        "&", "&amp;",
        "<", "&lt;",
        ">", "&gt;",
        `"`, "&#34;",
        "'", "&#39;",
    ).Replace(s)
}

// 将html部分转义标签还原为html特殊标签
// 参考：http://php.net/manual/zh/function.htmlspecialchars-decode.php
func SpecialCharsDecode(s string) string {
    return strings.NewReplacer(
        "&amp;", "&",
        "&lt;",  "<",
        "&gt;",  ">",
        "&#34;", `"`,
        "&#39;", "'",
    ).Replace(s)
}
