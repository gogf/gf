// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gogf/gf/v2/util/grand"
)

var (
	// octReg is the regular expression object for checks octal string.
	octReg = regexp.MustCompile(`\\[0-7]{3}`)
)

// Chr return the ascii string of a number(0-255).
//
// Example:
// Chr(65) -> "A"
func Chr(ascii int) string {
	return string([]byte{byte(ascii % 256)})
}

// Ord converts the first byte of a string to a value between 0 and 255.
//
// Example:
// Chr("A") -> 65
func Ord(char string) int {
	return int(char[0])
}

// OctStr converts string container octal string to its original string,
// for example, to Chinese string.
//
// Example:
// OctStr("\346\200\241") -> 怡
func OctStr(str string) string {
	return octReg.ReplaceAllStringFunc(
		str,
		func(s string) string {
			i, _ := strconv.ParseInt(s[1:], 8, 0)
			return string([]byte{byte(i)})
		},
	)
}

// Reverse returns a string which is the reverse of `str`.
//
// Example:
// Reverse("123456") -> "654321"
func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NumberFormat formats a number with grouped thousands.
// Parameter `decimals`: Sets the number of decimal points.
// Parameter `decPoint`: Sets the separator for the decimal point.
// Parameter `thousandsSep`: Sets the thousands' separator.
// See http://php.net/manual/en/function.number-format.php.
//
// Example:
// NumberFormat(1234.56, 2, ".", "")  -> 1234,56
// NumberFormat(1234.56, 2, ",", " ") -> 1 234,56
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

// Shuffle randomly shuffles a string.
// It considers parameter `str` as unicode string.
//
// Example:
// Shuffle("123456") -> "325164"
// Shuffle("123456") -> "231546"
// ...
func Shuffle(str string) string {
	runes := []rune(str)
	s := make([]rune, len(runes))
	for i, v := range grand.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// HideStr replaces part of the string `str` to `hide` by `percentage` from the `middle`.
// It considers parameter `str` as unicode string.
func HideStr(str string, percent int, hide string) string {
	// Handle email case
	var suffix string
	if idx := strings.IndexByte(str, '@'); idx >= 0 {
		suffix = str[idx:]
		str = str[:idx]
	}

	// Early return for edge cases
	if str == "" || percent <= 0 {
		return str + suffix
	}
	if percent >= 100 {
		return strings.Repeat(hide, len([]rune(str))) + suffix
	}

	rs := []rune(str)
	length := len(rs)
	if length == 0 {
		return str + suffix
	}

	// Calculate hideLen using the same logic as original (with floor)
	hideLen := (length * percent) / 100
	if hideLen == 0 {
		return str + suffix
	}

	// Calculate start position: mid - hideLen/2
	// This matches the original algorithm behavior
	mid := length / 2
	start := max(mid-hideLen/2, 0)

	end := start + hideLen
	if end > length {
		end = length
		start = max(length-hideLen, 0)
	}

	// Pre-calculate capacity to avoid reallocations
	var builder strings.Builder
	builder.Grow(len(str) + len(hide)*hideLen + len(suffix))

	// Build result string efficiently
	builder.WriteString(string(rs[:start]))
	if hide != "" {
		builder.WriteString(strings.Repeat(hide, hideLen))
	}
	builder.WriteString(string(rs[end:]))
	builder.WriteString(suffix)

	return builder.String()
}

// Nl2Br inserts HTML line breaks(`br`|<br />) before all newlines in a string:
// \n\r, \r\n, \r, \n.
// It considers parameter `str` as unicode string.
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
			if (i+1 < length) && ((v == r && runes[i+1] == n) || (v == n && runes[i+1] == r)) {
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

// WordWrap wraps a string to a given number of characters.
// This function supports cut parameters of both english and chinese punctuations.
// TODO: Enable custom cut parameter, see http://php.net/manual/en/function.wordwrap.php.
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
	for _, char := range str {
		switch {
		case char == '\n':
			if wordBuf.Len() == 0 {
				if current+spaceBuf.Len() > width {
					current = 0
				} else {
					current += spaceBuf.Len()
					_, _ = spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += spaceBuf.Len() + wordBuf.Len()
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			buf.WriteRune(char)
			current = 0

		case unicode.IsSpace(char):
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBuf.Len() + wordBuf.Len()
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			spaceBuf.WriteRune(char)

		case isPunctuation(char):
			wordBuf.WriteRune(char)
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBuf.Len() + wordBuf.Len()
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}

		default:
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
			_, _ = spaceBuf.WriteTo(buf)
		}
	} else {
		_, _ = spaceBuf.WriteTo(buf)
		_, _ = wordBuf.WriteTo(buf)
	}
	return buf.String()
}

func isPunctuation(char int32) bool {
	switch char {
	// English Punctuations.
	case ';', '.', ',', ':', '~':
		return true
	// Chinese Punctuations.
	case '；', '，', '。', '：', '？', '！', '…', '、':
		return true
	default:
		return false
	}
}
