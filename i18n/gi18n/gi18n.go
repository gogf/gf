// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gi18n implements internationalization and localization.
package gi18n

var (
	// defaultManager is the default i18n instance for package functions.
	defaultManager = Instance()
)

// SetPath sets the directory path storing i18n files.
func SetPath(path string) error {
	return defaultManager.SetPath(path)
}

// SetLanguage sets the language for translator.
func SetLanguage(language string) {
	defaultManager.SetLanguage(language)
}

// SetDelimiters sets the delimiters for translator.
func SetDelimiters(left, right string) {
	defaultManager.SetDelimiters(left, right)
}

// T is alias of Translate for convenience.
func T(content string, language ...string) string {
	return defaultManager.T(content, language...)
}

// Translate translates <content> with configured language.
// The parameter <language> specifies custom translation language ignoring configured language.
func Translate(content string, language ...string) string {
	return defaultManager.Translate(content, language...)
}

// GetValue retrieves and returns the configured content for given key and specified language.
// It returns an empty string if not found.
func GetContent(key string, language ...string) string {
	return defaultManager.GetContent(key, language...)
}
