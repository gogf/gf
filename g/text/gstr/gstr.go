// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gstr provides functions for string handling.
// 
// 字符串处理.
package gstr

import (
    "bytes"
    "fmt"
    "github.com/gogf/gf/g/util/grand"
    "math"
    "strconv"
    "strings"
    "unicode"
    "unicode/utf8"
)

// Replace returns a copy of the string <origin> with string <search> replaced by <replace>.
//
// 字符串替换(大小写敏感)
func Replace(origin, search, replace string, count...int) string {
    n := -1
    if len(count) > 0 {
        n = count[0]
    }
    return strings.Replace(origin, search, replace, n)
}

// Replace returns a copy of the string <origin> with string <search> replaced by <replace>
// with case-insensitive.
//
// 字符串替换(大小写不敏感)
func ReplaceI(origin, search, replace string, count...int) string {
    n := -1
    if len(count) > 0 {
        n = count[0]
    }
    if n == 0 {
        return origin
    }
    length      := len(search)
    searchLower := strings.ToLower(search)
    for {
        originLower := strings.ToLower(origin)
        if pos := strings.Index(originLower, searchLower); pos != -1 {
            origin = origin[ : pos] + replace + origin[pos + length : ]
            if n -= 1; n == 0 {
                break
            }
        } else {
            break
        }
    }
    return origin
}

// Count counts the number of <substr> appears in <s>. It returns 0 if no <substr> found in <s>.
//
// 计算字符串substr在字符串s中出现的次数，如果没有在s中找到substr，那么返回0。
func Count(s, substr string) int {
    return strings.Count(s, substr)
}

// Count counts the number of <substr> appears in <s>, case-insensitive. It returns 0 if no <substr> found in <s>.
//
// (非大小写敏感)计算字符串substr在字符串s中出现的次数，如果没有在s中找到substr，那么返回0。
func CountI(s, substr string) int {
    return strings.Count(ToLower(s), ToLower(substr))
}

// Replace string by array/slice.
//
// 使用map进行字符串替换(大小写敏感)
func ReplaceByArray(origin string, array []string) string {
    for i := 0; i < len(array); i += 2 {
        if i + 1 >= len(array) {
            break
        }
        origin = Replace(origin, array[i], array[i + 1])
    }
    return origin
}

// Replace string by array/slice with case-insensitive.
//
// 使用map进行字符串替换(大小写不敏感)
func ReplaceIByArray(origin string, array []string) string {
    for i := 0; i < len(array); i += 2 {
        if i + 1 >= len(array) {
            break
        }
        origin = ReplaceI(origin, array[i], array[i + 1])
    }
    return origin
}

// Replace string by map.
//
// 使用map进行字符串替换(大小写敏感)
func ReplaceByMap(origin string, replaces map[string]string) string {
    for k, v := range replaces {
        origin = Replace(origin, k, v)
    }
    return origin
}

// Replace string by map with case-insensitive.
//
// 使用map进行字符串替换(大小写不敏感)
func ReplaceIByMap(origin string, replaces map[string]string) string {
    for k, v := range replaces {
        origin = ReplaceI(origin, k, v)
    }
    return origin
}

// ToLower returns a copy of the string s with all Unicode letters mapped to their lower case.
// 字符串转换为小写
func ToLower(s string) string {
    return strings.ToLower(s)
}

// ToUpper returns a copy of the string s with all Unicode letters mapped to their upper case.
//
// 字符串转换为大写
func ToUpper(s string) string {
    return strings.ToUpper(s)
}

// UcFirst returns a copy of the string s with the first letter mapped to its upper case.
//
// 字符串首字母转换为大写
func UcFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    if IsLetterLower(s[0]) {
        return string(s[0] - 32) + s[1 :]
    }
    return s
}

// LcFirst returns a copy of the string s with the first letter mapped to its lower case.
//
// 字符串首字母转换为小写
func LcFirst(s string) string {
    if len(s) == 0 {
        return s
    }
    if IsLetterUpper(s[0]) {
        return string(s[0] + 32) + s[1 :]
    }
    return s
}

// Uppercase the first character of each word in a string.
//
// 大写字符串中每个单词的第一个字符。
func UcWords(str string) string {
    return strings.Title(str)
}

// IsLetterLower tests whether the given byte b is in lower case.
//
// 判断给定字符是否小写
func IsLetterLower(b byte) bool {
    if b >= byte('a') && b <= byte('z') {
        return true
    }
    return false
}

// IsLetterUpper tests whether the given byte b is in upper case.
//
// 判断给定字符是否大写
func IsLetterUpper(b byte) bool {
    if b >= byte('A') && b <= byte('Z') {
        return true
    }
    return false
}

// IsNumeric tests whether the given string s is numeric.
//
// 判断锁给字符串是否为数字.
func IsNumeric(s string) bool {
    length := len(s)
    if length == 0 {
        return false
    }
    for i := 0; i < len(s); i++ {
        if s[i] < byte('0') || s[i] > byte('9') {
            return false
        }
    }
    return true
}

// Returns the portion of string specified by the start and length parameters.
//
// 字符串截取，支持中文
func SubStr(str string, start int, length...int) (substr string) {
    // 将字符串的转换成[]rune
    rs  := []rune(str)
    lth := len(rs)
    // 简单的越界判断
    if start < 0 {
        start = 0
    }
    if start >= lth {
        start = lth
    }
    end := lth
    if len(length) > 0 {
        end = start + length[0]
        if end < start {
            end = lth
        }
    }
    if end > lth {
        end = lth
    }
    // 返回子串
    return string(rs[start : end])
}

// Returns the portion of string specified by the <length> parameters,
// if the length of str is greater than <length>,
// then the <suffix> will be appended to the result.
//
// 字符串长度截取限制，超过长度限制被截取并在字符串末尾追加指定的内容，支持中文
func StrLimit(str string, length int, suffix...string) (string) {
    rs := []rune(str)
    if len(str) < length {
        return str
    }
    addStr := "..."
    if len(suffix) > 0 {
        addStr = suffix[0]
    }
    return string(rs[0 : length]) + addStr
}

// Reverse a string.
//
// 字符串反转.
func Reverse(str string) string {
    runes := []rune(str)
    for i, j := 0, len(runes) - 1; i < j; i, j = i + 1, j - 1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// Format a number with grouped thousands.
// decimals: Sets the number of decimal points.
// decPoint: Sets the separator for the decimal point.
// thousandsSep: Sets the thousands separator.
// See http://php.net/manual/en/function.number-format.php.
//
// 以千位分隔符方式格式化一个数字.
func NumberFormat(number float64, decimals int, decPoint, thousandsSep string) string {
    neg := false
    if number < 0 {
        number = -number
        neg = true
    }
    // Will round off
    str := fmt.Sprintf("%."+strconv.Itoa(decimals)+"F", number)
    prefix, suffix := "", ""
    if decimals > 0 {
        prefix = str[ : len(str) - (decimals + 1)]
        suffix = str[len(str) - decimals : ]
    } else {
        prefix = str
    }
    sep := []byte(thousandsSep)
    n, l1, l2 := 0, len(prefix), len(sep)
    // thousands sep num
    c := (l1 - 1) / 3
    tmp := make([]byte, l2*c+l1)
    pos := len(tmp) - 1
    for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
        if l2 > 0 && n > 0 && n%3 == 0 {
            for j := range sep {
                tmp[pos] = sep[l2-j-1]
                pos--
            }
        }
        tmp[pos] = prefix[i]
    }
    s := string(tmp)
    if decimals > 0 {
        s += decPoint + suffix
    }
    if neg {
        s = "-" + s
    }

    return s
}

// Split a string into smaller chunks.
// Can be used to split a string into smaller chunks which is useful for
// e.g. converting BASE64 string output to match RFC 2045 semantics.
// It inserts end every chunkLen characters.
//
// 将字符串分割成小块。使用此函数将字符串分割成小块非常有用。
// 例如将BASE64的输出转换成符合RFC2045语义的字符串。
// 它会在每 chunkLen 个字符后边插入 end。
func ChunkSplit(body string, chunkLen int, end string) string {
    if end == "" {
        end = "\r\n"
    }
    runes, endRunes := []rune(body), []rune(end)
    l := len(runes)
    if l <= 1 || l < chunkLen {
        return body + end
    }
    ns := make([]rune, 0, len(runes) + len(endRunes))
    for i := 0; i < l; i += chunkLen {
        if i + chunkLen > l {
            ns = append(ns, runes[i : ]...)
        } else {
            ns = append(ns, runes[i : i + chunkLen]...)
        }
        ns = append(ns, endRunes...)
    }
    return string(ns)
}

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
//
// 比较两个字符串。
func Compare(a, b string) int {
    return strings.Compare(a, b)
}

// Equal reports whether s and t, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, case-insensitive.
//
// 比较两个字符串是否相等（不区分大小写）。
func Equal(a, b string) bool {
    return strings.EqualFold(a, b)
}

// Return the words used in a string.
//
// 分割字符串中的单词。
func Fields(str string) []string {
    return strings.Fields(str)
}

// Contains reports whether substr is within str.
//
// 判断是否substr存在于str中。
func Contains(str, substr string) bool {
    return strings.Contains(str, substr)
}

// Contains reports whether substr is within str, case-insensitive.
//
// 判断是否substr存在于str中(不区分大小写)。
func ContainsI(str, substr string) bool {
    return PosI(str, substr) != -1
}

// ContainsAny reports whether any Unicode code points in chars are within s.
//
// 判断是否s中是否包含chars指定的任意字符。
func ContainsAny(s, chars string) bool {
    return strings.ContainsAny(s, chars)
}

// Return information about words used in a string.
//
// 返回字符串中单词的使用情况。
func CountWords(str string) map[string]int {
    m      := make(map[string]int)
    buffer := bytes.NewBuffer(nil)
    for _, rune := range []rune(str) {
        if unicode.IsSpace(rune) {
            if buffer.Len() > 0 {
                m[buffer.String()]++
                buffer.Reset()
            }
        } else {
            buffer.WriteRune(rune)
        }
    }
    if buffer.Len() > 0 {
        m[buffer.String()]++
    }
    return m
}

// Return information about words used in a string.
//
// 返回字符串中字符的使用情况。
func CountChars(str string, noSpace...bool) map[string]int {
    m := make(map[string]int)
    countSpace := true
    if len(noSpace) > 0 && noSpace[0] {
        countSpace = false
    }
    for _, rune := range []rune(str) {
        if !countSpace && unicode.IsSpace(rune) {
            continue
        }
        m[string(rune)]++
    }
    return m
}

// Wraps a string to a given number of characters.
// TODO: Enable cut param, see http://php.net/manual/en/function.wordwrap.php.
//
// 使用字符串断点将字符串打断为指定数量的字串。
func WordWrap(str string, width int, br string) string {
    if br == "" {
        br = "\n"
    }
    init := make([]byte, 0, len(str))
    buf := bytes.NewBuffer(init)
    var current int
    var wordBuf, spaceBuf bytes.Buffer
    for _, char := range str {
        if char == '\n' {
            if wordBuf.Len() == 0 {
                if current + spaceBuf.Len() > width {
                    current = 0
                } else {
                    current += spaceBuf.Len()
                    spaceBuf.WriteTo(buf)
                }
                spaceBuf.Reset()
            } else {
                current += spaceBuf.Len() + wordBuf.Len()
                spaceBuf.WriteTo(buf)
                spaceBuf.Reset()
                wordBuf.WriteTo(buf)
                wordBuf.Reset()
            }
            buf.WriteRune(char)
            current = 0
        } else if unicode.IsSpace(char) {
            if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
                current += spaceBuf.Len() + wordBuf.Len()
                spaceBuf.WriteTo(buf)
                spaceBuf.Reset()
                wordBuf.WriteTo(buf)
                wordBuf.Reset()
            }
            spaceBuf.WriteRune(char)
        } else {
            wordBuf.WriteRune(char)
            if current + spaceBuf.Len()+wordBuf.Len() > width && wordBuf.Len() < width {
                buf.WriteString(br)
                current = 0
                spaceBuf.Reset()
            }
        }
    }

    if wordBuf.Len() == 0 {
        if current + spaceBuf.Len() <= width {
            spaceBuf.WriteTo(buf)
        }
    } else {
        spaceBuf.WriteTo(buf)
        wordBuf.WriteTo(buf)
    }
    return buf.String()
}

// Get string length of unicode.
//
// UTF-8字符串长度。
func RuneLen(str string) int {
    return utf8.RuneCountInString(str)
}

// Repeat returns a new string consisting of multiplier copies of the string input.
//
// 按照指定大小创建重复的字符串。
func Repeat(input string, multiplier int) string {
    return strings.Repeat(input, multiplier)
}

// Returns part of haystack string starting from and including the first occurrence of needle to the end of haystack.
// See http://php.net/manual/en/function.strstr.php.
//
// 查找字符串的首次出现。返回 haystack 字符串从 needle 第一次出现的位置开始到 haystack 结尾的字符串。
func Str(haystack string, needle string) string {
    if needle == "" {
        return ""
    }
    idx := strings.Index(haystack, needle)
    if idx == -1 {
        return ""
    }
    return haystack[idx + len([]byte(needle)) - 1 : ]
}

// Randomly shuffles a string.
//
// 将字符串打乱。
func Shuffle(str string) string {
    runes := []rune(str)
    s     := make([]rune, len(runes))
    for i, v := range grand.Perm(len(runes)) {
        s[i] = runes[v]
    }
    return string(s)
}

// Split a string by a string, to an array.
//
// 此函数返回由字符串组成的数组，每个元素都是 str 的一个子串，它们被字符串 delimiter 作为边界点分割出来。
func Split(str, delimiter string) []string {
    return strings.Split(str, delimiter)
}

// Join concatenates the elements of a to create a single string. The separator string
// sep is placed between elements in the resulting string.
//
// 用sep将字符串数组array连接为一个字符串。
func Join(array []string, sep string) string {
    return strings.Join(array, sep)
}

// Split a string by a string, to an array.
// See http://php.net/manual/en/function.explode.php.
//
// 此函数返回由字符串组成的数组，每个元素都是 str 的一个子串，它们被字符串 delimiter 作为边界点分割出来。
func Explode(delimiter, str string) []string {
    return Split(str, delimiter)
}

// Join array elements with a string.
// http://php.net/manual/en/function.implode.php
//
// 用glue将字符串数组pieces连接为一个字符串。
func Implode(glue string, pieces []string) string {
    return strings.Join(pieces, glue)
}

// Generate a single-byte string from a number.
//
// 返回相对应于 ascii 所指定的单个字符。
func Chr(ascii int) string {
    return string(ascii)
}

// Convert the first byte of a string to a value between 0 and 255.
//
// 解析 char 二进制值第一个字节为 0 到 255 范围的无符号整型类型。
func Ord(char string) int {
    return int(char[0])
}

// HideStr replaces part of the the string by percentage from the middle.
//
// 按照百分比从字符串中间向两边隐藏字符(主要用于姓名、手机号、邮箱地址、身份证号等的隐藏)，支持utf-8中文，支持email格式。
func HideStr(str string, percent int, hide string) string {
    array := strings.Split(str, "@")
    if len(array) > 1 {
        str = array[0]
    }
    rs       := []rune(str)
    length   := len(rs)
    mid      := math.Floor(float64(length/2))
    hideLen  := int(math.Floor(float64(length) * (float64(percent)/100)))
    start    := int(mid - math.Floor(float64(hideLen) / 2))
    hideStr  := []rune("")
    hideRune := []rune(hide)
    for i := 0; i < int(hideLen); i++ {
        hideStr = append(hideStr, hideRune...)
    }
    buffer := bytes.NewBuffer(nil)
    buffer.WriteString(string(rs[0 : start]))
    buffer.WriteString(string(hideStr))
    buffer.WriteString(string(rs[start + hideLen : ]))
    if len(array) > 1 {
        buffer.WriteString("@" + array[1])
    }
    return buffer.String()
}

// Inserts HTML line breaks before all newlines in a string.
// \n\r, \r\n, \r, \n
//
// 在字符串 string 所有新行之前插入 '<br />' 或 '<br>'，并返回。
func Nl2Br(str string, isXhtml...bool) string {
    r, n, runes := '\r', '\n', []rune(str)
    var br []byte
    if len(isXhtml) > 0 && isXhtml[0] {
        br = []byte("<br />")
    } else {
        br = []byte("<br>")
    }
    skip := false
    length := len(runes)
    var buf bytes.Buffer
    for i, v := range runes {
        if skip {
            skip = false
            continue
        }
        switch v {
            case n, r:
                if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
                    buf.Write(br)
                    skip = true
                    continue
                }
                buf.Write(br)
            default:
                buf.WriteRune(v)
        }
    }
    return buf.String()
}

// Quote string with slashes.
//
// 转义字符串中的单引号（'）、双引号（"）、反斜线（\）与 NUL（NULL 字符）。
func AddSlashes(str string) string {
    var buf bytes.Buffer
    for _, char := range str {
        switch char {
            case '\'', '"', '\\':
                buf.WriteRune('\\')
        }
        buf.WriteRune(char)
    }
    return buf.String()
}

// Un-quotes a quoted string.
//
// 反转义字符串。
func StripSlashes(str string) string {
    var buf bytes.Buffer
    l, skip := len(str), false
    for i, char := range str {
        if skip {
            skip = false
        } else if char == '\\' {
            if i + 1 < l && str[i + 1] == '\\' {
                skip = true
            }
            continue
        }
        buf.WriteRune(char)
    }
    return buf.String()
}

// Returns a version of str with a backslash character (\) before every character that is among:
// .\+*?[^]($)
//
// 转义字符串，转义的特殊字符包括：.\+*?[^]($)。
func QuoteMeta(str string) string {
    var buf bytes.Buffer
    for _, char := range str {
        switch char {
        case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
            buf.WriteRune('\\')
        }
        buf.WriteRune(char)
    }
    return buf.String()
}