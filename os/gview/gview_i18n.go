// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import "github.com/gogf/gf/util/gconv"

// i18nTranslate translate the content with i18n feature.
func (view *View) i18nTranslate(content string, params Params) string {
	if view.config.I18nManager != nil {
		if v, ok := params["I18nLanguage"]; ok {
			language := gconv.String(v)
			if language != "" {
				return view.config.I18nManager.T(content, language)
			}
		}
		return view.config.I18nManager.T(content)
	}
	return content
}
