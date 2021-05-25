// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import (
	"context"
	"github.com/gogf/gf/i18n/gi18n"
	"github.com/gogf/gf/util/gconv"
)

const (
	i18nLanguageVariableName = "I18nLanguage"
)

// i18nTranslate translate the content with i18n feature.
func (view *View) i18nTranslate(ctx context.Context, content string, variables Params) string {
	if view.config.I18nManager != nil {
		// Compatible with old version.
		if language, ok := variables[i18nLanguageVariableName]; ok {
			ctx = gi18n.WithLanguage(ctx, gconv.String(language))
		}
		return view.config.I18nManager.T(ctx, content)
	}
	return content
}

// setI18nLanguageFromCtx retrieves language name from context and sets it to template variables map.
func (view *View) setI18nLanguageFromCtx(ctx context.Context, variables map[string]interface{}) {
	if language, ok := variables[i18nLanguageVariableName]; !ok {
		if language = gi18n.LanguageFromCtx(ctx); language != "" {
			variables[i18nLanguageVariableName] = language
		}
	}
}
