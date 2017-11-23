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
