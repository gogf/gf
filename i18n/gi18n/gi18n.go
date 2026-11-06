// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gi18n implements internationalization and localization.
package gi18n

import "context"

// SetPath sets the directory path storing i18n files.
func SetPath(path string) error {
	return Instance().SetPath(path)
}

// SetLanguage sets the language for translator.
func SetLanguage(language string) {
	Instance().SetLanguage(language)
}

// SetDelimiters sets the delimiters for translator.
func SetDelimiters(left, right string) {
	Instance().SetDelimiters(left, right)
}

// T is alias of Translate for convenience.
func T(ctx context.Context, content string) string {
	return Instance().T(ctx, content)
}

// Tf is alias of TranslateFormat for convenience.
func Tf(ctx context.Context, format string, values ...any) string {
	return Instance().TranslateFormat(ctx, format, values...)
}

// TranslateFormat translates, formats and returns the `format` with configured language
// and given `values`.
func TranslateFormat(ctx context.Context, format string, values ...any) string {
	return Instance().TranslateFormat(ctx, format, values...)
}

// Tm is alias of TranslateMap for convenience.
func Tm(ctx context.Context, content string, valMap map[string]any) string {
	return Instance().TranslateMap(ctx, content, valMap)
}

// TranslateMap parameter replacement is order-independent.
// replaced param format is {$key}.  for example: content is "my name is {$name}.",the valMap is map[string]interface{} {"name":"alias"}
// the result of translation is "my name is alias."
func TranslateMap(ctx context.Context, content string, valMap map[string]any) string {
	return Instance().TranslateMap(ctx, content, valMap)
}

// Translate translates `content` with configured language and returns the translated content.
func Translate(ctx context.Context, content string) string {
	return Instance().Translate(ctx, content)
}

// GetContent retrieves and returns the configured content for given key and specified language.
// It returns an empty string if not found.
func GetContent(ctx context.Context, key string) string {
	return Instance().GetContent(ctx, key)
}
