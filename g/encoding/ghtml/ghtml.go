// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// HTML编码
package ghtml

import "strings"

// 将html中的特殊标签转换为html转义标签
func SpecialChars(s string) string {
    return strings.NewReplacer(
        "&", "&amp;",
        "<", "&lt;",
        ">", "&gt;",
        `"`, "&#34;",
        "'", "&#39;",
    ).Replace(s)
}

// 将html转义标签还原为html特殊标签
func SpecialCharsDecode(s string) string {
    return strings.NewReplacer(
        "&amp;", "&",
        "&lt;",  "<",
        "&gt;",  ">",
        "&#34;", `"`,
        "&#39;", "'",
    ).Replace(s)
}
