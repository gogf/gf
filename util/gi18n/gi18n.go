// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gi18n implements internationalization and localization.
package gi18n

var (
	defaultTranslator = Instance()
)

// SetPath sets the directory path storing i18n files.
func SetPath(path string) error {
	return defaultTranslator.SetPath(path)
}

// SetLanguage sets the language for translator.
func SetLanguage(language string) {
	defaultTranslator.SetLanguage(language)
}

// SetDelimiters sets the delimiters for translator.
func SetDelimiters(left, right string) {
	defaultTranslator.SetDelimiters(left, right)
}

// T is alias of Translate.
func T(content string, language ...string) string {
	return defaultTranslator.T(content, language...)
}

// Translate translates <content> with configured language.
// The parameter <language> specifies custom translation language ignoring configured language.
func Translate(content string, language ...string) string {
	return defaultTranslator.Translate(content, language...)
}
