// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"
)

var (
	// defaultTrimChars are the characters which are stripped by Trim* functions in default.
	defaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})
)

// Trim strips whitespace (or other characters) from the beginning and end of a string.
// The optional parameter <characterMask> specifies the additional stripped characters.
func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.Trim(str, defaultTrimChars)
	} else {
		return strings.Trim(str, defaultTrimChars+characterMask[0])
	}
}

// TrimStr strips all of the given <cut> string from the beginning and end of a string.
// Note that it does not strips the whitespaces of its beginning or end.
func TrimStr(str string, cut string) string {
	return TrimLeftStr(TrimRightStr(str, cut), cut)
}

// TrimLeft strips whitespace (or other characters) from the beginning of a string.
func TrimLeft(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimLeft(str, defaultTrimChars)
	} else {
		return strings.TrimLeft(str, defaultTrimChars+characterMask[0])
	}
}

// TrimLeftStr strips all of the given <cut> string from the beginning of a string.
// Note that it does not strips the whitespaces of its beginning.
func TrimLeftStr(str string, cut string) string {
	var lenCut = len(cut)
	for len(str) >= lenCut && str[0:lenCut] == cut {
		str = str[lenCut:]
	}
	return str
}

// TrimRight strips whitespace (or other characters) from the end of a string.
func TrimRight(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.TrimRight(str, defaultTrimChars)
	} else {
		return strings.TrimRight(str, defaultTrimChars+characterMask[0])
	}
}

// TrimRightStr strips all of the given <cut> string from the end of a string.
// Note that it does not strips the whitespaces of its end.
func TrimRightStr(str string, cut string) string {
	var lenStr = len(str)
	var lenCut = len(cut)
	for lenStr >= lenCut && str[lenStr-lenCut:lenStr] == cut {
		lenStr = lenStr - lenCut
		str = str[:lenStr]

	}
	return str
}
