package gutil

import (
    "time"
    "math/rand"
    "strings"
)

// 框架自定义函数库

// 获得一个 min, max 之间的随机数
func Rand (min, max int) int {
    //fmt.Printf("min: %d, max: %d\n", min, max)
    if min >= max {
        return min
    }
    rand.Seed(time.Now().UnixNano())
    n := rand.Intn(max)
    if n < min {
        return Rand(min, max)
    }
    return n
}

// 将html中的特殊标签转换为html转义标签
func HtmlSpecialChars(s string) string {
    return strings.NewReplacer(
        "&", "&amp;",
        "<", "&lt;",
        ">", "&gt;",
        `"`, "&#34;",
        "'", "&#39;",
    ).Replace(s)
}

// 将html转义标签还原为html特殊标签
func HtmlSpecialCharsDecode(s string) string {
    return strings.NewReplacer(
        "&amp;", "&",
        "&lt;",  "<",
        "&gt;",  ">",
        "&#34;", `"`,
        "&#39;", "'",
    ).Replace(s)
}
