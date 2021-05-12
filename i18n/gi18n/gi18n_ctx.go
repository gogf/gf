// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gi18n implements internationalization and localization.
package gi18n

import "context"

const (
	ctxLanguage = "I18nLanguage"
)

// WithLanguage append language setting to the context and returns a new context.
func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, ctxLanguage, language)
}

// LanguageFromCtx retrieves and returns language name from context.
// It returns an empty string if it is not set previously.
func LanguageFromCtx(ctx context.Context) string {
	v := ctx.Value(ctxLanguage)
	if v != nil {
		return v.(string)
	}
	return ""
}
