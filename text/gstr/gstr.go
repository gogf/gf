// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gstr provides functions for string handling.
package gstr

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gogf/gf/internal/utils"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/util/grand"
)

// Replace returns a copy of the string <origin>
// in which string <search> replaced by <replace> case-sensitively.
func Replace(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	return strings.Replace(origin, search, replace, n)
}

// Replace returns a copy of the string <origin>
// in which string <search> replaced by <replace> case-insensitively.
func ReplaceI(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	if n == 0 {
		return origin
	}
	length := len(search)
	searchLower := strings.ToLower(search)
	for {
		originLower := strings.ToLower(origin)
		if pos := strings.Index(originLower, searchLower); pos != -1 {
			origin = origin[:pos] + replace + origin[pos+length:]
			if n--; n == 0 {
				break
			}
		} else {
			break
		}
	}
	return origin
}

// Count counts the number of <substr> appears in <s>.
// It returns 0 if no <substr> found in <s>.
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// CountI counts the number of <substr> appears in <s>, case-insensitively.
// It returns 0 if no <substr> found in <s>.
func CountI(s, substr string) int {
	return strings.Count(ToLower(s), ToLower(substr))
}

// ReplaceByArray returns a copy of <origin>,
// which is replaced by a slice in order, case-sensitively.
func ReplaceByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = Replace(origin, array[i], array[i+1])
	}
	return origin
}

// ReplaceIByArray returns a copy of <origin>,
// which is replaced by a slice in order, case-insensitively.
func ReplaceIByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = ReplaceI(origin, array[i], array[i+1])
	}
	return origin
}

// ReplaceByMap returns a copy of <origin>,
// which is replaced by a map in unordered way, case-sensitively.
func ReplaceByMap(origin string, replaces map[string]string) string {
	return utils.ReplaceByMap(origin, replaces)
}

// ReplaceIByMap returns a copy of <origin>,
// which is replaced by a map in unordered way, case-insensitively.
func ReplaceIByMap(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = ReplaceI(origin, k, v)
	}
	return origin
}

// ToLower returns a copy of the string s with all Unicode letters mapped to their lower case.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper returns a copy of the string s with all Unicode letters mapped to their upper case.
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// UcFirst returns a copy of the string s with the first letter mapped to its upper case.
func UcFirst(s string) string {
	return utils.UcFirst(s)
}

// LcFirst returns a copy of the string s with the first letter mapped to its lower case.
func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if IsLetterUpper(s[0]) {
		return string(s[0]+32) + s[1:]
	}
	return s
}

// UcWords uppercase the first character of each word in a string.
func UcWords(str string) string {
	return strings.Title(str)
}

// IsLetterLower tests whether the given byte b is in lower case.
func IsLetterLower(b byte) bool {
	return utils.IsLetterLower(b)
}

// IsLetterUpper tests whether the given byte b is in upper case.
func IsLetterUpper(b byte) bool {
	return utils.IsLetterUpper(b)
}

// IsNumeric tests whether the given string s is numeric.
func IsNumeric(s string) bool {
	return utils.IsNumeric(s)
}

// SubStr returns a portion of string <str> specified by the <start> and <length> parameters.
func SubStr(str string, start int, length ...int) (substr string) {
	lth := len(str)

	// Simple border checks.
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
	return str[start:end]
}

// SubStrRune returns a portion of string <str> specified by the <start> and <length> parameters.
// SubStrRune considers parameter <str> as unicode string.
func SubStrRune(str string, start int, length ...int) (substr string) {
	// Converting to []rune to support unicode.
	rs := []rune(str)
	lth := len(rs)

	// Simple border checks.
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
	return string(rs[start:end])
}

// StrLimit returns a portion of string <str> specified by <length> parameters, if the length
// of <str> is greater than <length>, then the <suffix> will be appended to the result string.
func StrLimit(str string, length int, suffix ...string) string {
	if len(str) < length {
		return str
	}
	addStr := "..."
	if len(suffix) > 0 {
		addStr = suffix[0]
	}
	return str[0:length] + addStr
}

// StrLimitRune returns a portion of string <str> specified by <length> parameters, if the length
// of <str> is greater than <length>, then the <suffix> will be appended to the result string.
// StrLimitRune considers parameter <str> as unicode string.
func StrLimitRune(str string, length int, suffix ...string) string {
	rs := []rune(str)
	if len(rs) < length {
		return str
	}
	addStr := "..."
	if len(suffix) > 0 {
		addStr = suffix[0]
	}
	return string(rs[0:length]) + addStr
}

// Reverse returns a string which is the reverse of <str>.
func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NumberFormat formats a number with grouped thousands.
// <decimals>: Sets the number of decimal points.
// <decPoint>: Sets the separator for the decimal point.
// <thousandsSep>: Sets the thousands separator.
// See http://php.net/manual/en/function.number-format.php.
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
		prefix = str[:len(str)-(decimals+1)]
		suffix = str[len(str)-decimals:]
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

// ChunkSplit splits a string into smaller chunks.
// Can be used to split a string into smaller chunks which is useful for
// e.g. converting BASE64 string output to match RFC 2045 semantics.
// It inserts end every chunkLen characters.
// It considers parameter <body> and <end> as unicode string.
func ChunkSplit(body string, chunkLen int, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, endRunes := []rune(body), []rune(end)
	l := len(runes)
	if l <= 1 || l < chunkLen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(endRunes))
	for i := 0; i < l; i += chunkLen {
		if i+chunkLen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunkLen]...)
		}
		ns = append(ns, endRunes...)
	}
	return string(ns)
}

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a, b string) int {
	return strings.Compare(a, b)
}

// Equal reports whether <a> and <b>, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, case-insensitively.
func Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}

// Fields returns the words used in a string as slice.
func Fields(str string) []string {
	return strings.Fields(str)
}

// HasPrefix tests whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix tests whether the string s ends with suffix.
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// CountWords returns information about words' count used in a string.
// It considers parameter <str> as unicode string.
func CountWords(str string) map[string]int {
	m := make(map[string]int)
	buffer := bytes.NewBuffer(nil)
	for _, r := range []rune(str) {
		if unicode.IsSpace(r) {
			if buffer.Len() > 0 {
				m[buffer.String()]++
				buffer.Reset()
			}
		} else {
			buffer.WriteRune(r)
		}
	}
	if buffer.Len() > 0 {
		m[buffer.String()]++
	}
	return m
}

// CountChars returns information about chars' count used in a string.
// It considers parameter <str> as unicode string.
func CountChars(str string, noSpace ...bool) map[string]int {
	m := make(map[string]int)
	countSpace := true
	if len(noSpace) > 0 && noSpace[0] {
		countSpace = false
	}
	for _, r := range []rune(str) {
		if !countSpace && unicode.IsSpace(r) {
			continue
		}
		m[string(r)]++
	}
	return m
}

// WordWrap wraps a string to a given number of characters.
// TODO: Enable cut parameter, see http://php.net/manual/en/function.wordwrap.php.
func WordWrap(str string, width int, br string) string {
	if br == "" {
		br = "\n"
	}
	var (
		current           int
		wordBuf, spaceBuf bytes.Buffer
		init              = make([]byte, 0, len(str))
		buf               = bytes.NewBuffer(init)
	)
	for _, char := range []rune(str) {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+spaceBuf.Len() > width {
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
			if current+spaceBuf.Len()+wordBuf.Len() > width && wordBuf.Len() < width {
				buf.WriteString(br)
				current = 0
				spaceBuf.Reset()
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+spaceBuf.Len() <= width {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}
	return buf.String()
}

// RuneLen returns string length of unicode.
// Deprecated, use LenRune instead.
func RuneLen(str string) int {
	return LenRune(str)
}

// LenRune returns string length of unicode.
func LenRune(str string) int {
	return utf8.RuneCountInString(str)
}

// Repeat returns a new string consisting of multiplier copies of the string input.
func Repeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// Str returns part of <haystack> string starting from and including
// the first occurrence of <needle> to the end of <haystack>.
// See http://php.net/manual/en/function.strstr.php.
func Str(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}

// Shuffle randomly shuffles a string.
// It considers parameter <str> as unicode string.
func Shuffle(str string) string {
	runes := []rune(str)
	s := make([]rune, len(runes))
	for i, v := range grand.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// Split splits string <str> by a string <delimiter>, to an array.
func Split(str, delimiter string) []string {
	return strings.Split(str, delimiter)
}

// SplitAndTrim splits string <str> by a string <delimiter> to an array,
// and calls Trim to every element of this array. It ignores the elements
// which are empty after Trim.
func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = Trim(v, characterMask...)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

// SplitAndTrimSpace splits string <str> by a string <delimiter> to an array,
// and calls TrimSpace to every element of this array.
// Deprecated.
func SplitAndTrimSpace(str, delimiter string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = strings.TrimSpace(v)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

// Join concatenates the elements of <array> to create a single string. The separator string
// <sep> is placed between elements in the resulting string.
func Join(array []string, sep string) string {
	return strings.Join(array, sep)
}

// JoinAny concatenates the elements of <array> to create a single string. The separator string
// <sep> is placed between elements in the resulting string.
//
// The parameter <array> can be any type of slice, which be converted to string array.
func JoinAny(array interface{}, sep string) string {
	return strings.Join(gconv.Strings(array), sep)
}

// Explode splits string <str> by a string <delimiter>, to an array.
// See http://php.net/manual/en/function.explode.php.
func Explode(delimiter, str string) []string {
	return Split(str, delimiter)
}

// Implode joins array elements <pieces> with a string <glue>.
// http://php.net/manual/en/function.implode.php
func Implode(glue string, pieces []string) string {
	return strings.Join(pieces, glue)
}

// Chr return the ascii string of a number(0-255).
func Chr(ascii int) string {
	return string([]byte{byte(ascii % 256)})
}

// Ord converts the first byte of a string to a value between 0 and 255.
func Ord(char string) int {
	return int(char[0])
}

// HideStr replaces part of the the string <str> to <hide> by <percentage> from the <middle>.
// It considers parameter <str> as unicode string.
func HideStr(str string, percent int, hide string) string {
	array := strings.Split(str, "@")
	if len(array) > 1 {
		str = array[0]
	}
	var (
		rs       = []rune(str)
		length   = len(rs)
		mid      = math.Floor(float64(length / 2))
		hideLen  = int(math.Floor(float64(length) * (float64(percent) / 100)))
		start    = int(mid - math.Floor(float64(hideLen)/2))
		hideStr  = []rune("")
		hideRune = []rune(hide)
	)
	for i := 0; i < hideLen; i++ {
		hideStr = append(hideStr, hideRune...)
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(string(rs[0:start]))
	buffer.WriteString(string(hideStr))
	buffer.WriteString(string(rs[start+hideLen:]))
	if len(array) > 1 {
		buffer.WriteString("@" + array[1])
	}
	return buffer.String()
}

// Nl2Br inserts HTML line breaks(<br>|<br />) before all newlines in a string:
// \n\r, \r\n, \r, \n.
// It considers parameter <str> as unicode string.
func Nl2Br(str string, isXhtml ...bool) string {
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

// AddSlashes quotes chars('"\) with slashes.
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

// StripSlashes un-quotes a quoted string by AddSlashes.
func StripSlashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// QuoteMeta returns a version of str with a backslash character (\)
// before every character that is among: .\+*?[^]($)
func QuoteMeta(str string, chars ...string) string {
	var buf bytes.Buffer
	for _, char := range str {
		if len(chars) > 0 {
			for _, c := range chars[0] {
				if c == char {
					buf.WriteRune('\\')
					break
				}
			}
		} else {
			switch char {
			case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
				buf.WriteRune('\\')
			}
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// SearchArray searches string <s> in string slice <a> case-sensitively,
// returns its index in <a>.
// If <s> is not found in <a>, it returns -1.
func SearchArray(a []string, s string) int {
	for i, v := range a {
		if s == v {
			return i
		}
	}
	return -1
}

// InArray checks whether string <s> in slice <a>.
func InArray(a []string, s string) bool {
	return SearchArray(a, s) != -1
}
