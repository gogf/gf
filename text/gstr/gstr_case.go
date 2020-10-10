// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
//   | Function                          | Result             |
//   |-----------------------------------|--------------------|
//   | SnakeCase(s)                      | any_kind_of_string |
//   | SnakeScreamingCase(s)             | ANY_KIND_OF_STRING |
//   | KebabCase(s)                      | any-kind-of-string |
//   | KebabScreamingCase(s)             | ANY-KIND-OF-STRING |
//   | DelimitedCase(s, '.')             | any.kind.of.string |
//   | DelimitedScreamingCase(s, '.')    | ANY.KIND.OF.STRING |
//   | CamelCase(s)                      | AnyKindOfString    |
//   | CamelLowerCase(s)                 | anyKindOfString    |
//   | SnakeFirstUpperCase(RGBCodeMd5)   | rgb_code_md5       |

package gstr

import (
	"regexp"
	"strings"
)

var (
	numberSequence      = regexp.MustCompile(`([a-zA-Z]{0,1})(\d+)([a-zA-Z]{0,1})`)
	firstCamelCaseStart = regexp.MustCompile(`([A-Z]+)([A-Z]?[_a-z\d]+)|$`)
	firstCamelCaseEnd   = regexp.MustCompile(`([\w\W]*?)([_]?[A-Z]+)$`)
)

// CamelCase converts a string to CamelCase.
func CamelCase(s string) string {
	return toCamelInitCase(s, true)
}

// CamelLowerCase converts a string to lowerCamelCase.
func CamelLowerCase(s string) string {
	if s == "" {
		return s
	}
	if r := rune(s[0]); r >= 'A' && r <= 'Z' {
		s = strings.ToLower(string(r)) + s[1:]
	}
	return toCamelInitCase(s, false)
}

// SnakeCase converts a string to snake_case.
func SnakeCase(s string) string {
	return DelimitedCase(s, '_')
}

// SnakeScreamingCase converts a string to SNAKE_CASE_SCREAMING.
func SnakeScreamingCase(s string) string {
	return DelimitedScreamingCase(s, '_', true)
}

// SnakeFirstUpperCase converts a string from RGBCodeMd5 to rgb_code_md5.
// The length of word should not be too long
// TODO for efficiency should change regexp to traversing string in future
func SnakeFirstUpperCase(word string, underscore ...string) string {
	replace := "_"
	if len(underscore) > 0 {
		replace = underscore[0]
	}

	m := firstCamelCaseEnd.FindAllStringSubmatch(word, 1)
	if len(m) > 0 {
		word = m[0][1] + replace + TrimLeft(ToLower(m[0][2]), replace)
	}

	for {
		m := firstCamelCaseStart.FindAllStringSubmatch(word, 1)
		if len(m) > 0 && m[0][1] != "" {
			w := strings.ToLower(m[0][1])
			w = string(w[:len(w)-1]) + replace + string(w[len(w)-1])

			word = strings.Replace(word, m[0][1], w, 1)
		} else {
			break
		}
	}

	return TrimLeft(word, replace)
}

// KebabCase converts a string to kebab-case
func KebabCase(s string) string {
	return DelimitedCase(s, '-')
}

// KebabScreamingCase converts a string to KEBAB-CASE-SCREAMING.
func KebabScreamingCase(s string) string {
	return DelimitedScreamingCase(s, '-', true)
}

// DelimitedCase converts a string to snake.case.delimited.
func DelimitedCase(s string, del uint8) string {
	return DelimitedScreamingCase(s, del, false)
}

// DelimitedScreamingCase converts a string to DELIMITED.SCREAMING.CASE or delimited.screaming.case.
func DelimitedScreamingCase(s string, del uint8, screaming bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			if (v >= 'A' && v <= 'Z' && next >= 'a' && next <= 'z') || (v >= 'a' && v <= 'z' && next >= 'A' && next <= 'Z') {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && n[len(n)-1] != del && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				n += string(del) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(del)
			}
		} else if v == ' ' || v == '_' || v == '-' || v == '.' {
			// replace spaces/underscores with delimiters
			n += string(del)
		} else {
			n = n + string(v)
		}
	}

	if screaming {
		n = strings.ToUpper(n)
	} else {
		n = strings.ToLower(n)
	}
	return n
}

func addWordBoundariesToNumbers(s string) string {
	r := numberSequence.ReplaceAllFunc([]byte(s), func(bytes []byte) []byte {
		var result []byte
		match := numberSequence.FindSubmatch(bytes)
		if len(match[1]) > 0 {
			result = append(result, match[1]...)
			result = append(result, []byte(" ")...)
		}
		result = append(result, match[2]...)
		if len(match[3]) > 0 {
			result = append(result, []byte(" ")...)
			result = append(result, match[3]...)
		}
		return result
	})
	return string(r)
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' || v == '.' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}
