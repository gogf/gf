// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gpage"
)

// GetPage creates and returns the pagination object for given <totalSize> and <pageSize>.
// NOTE THAT the page parameter name from client is constantly defined as gpage.DefaultPageName
// for simplification and convenience.
func (r *Request) GetPage(totalSize, pageSize int) *gpage.Page {
	// It must have Router object attribute.
	if r.Router == nil {
		panic("Router object not found")
	}
	var (
		url            = *r.URL
		urlTemplate    = url.Path
		uriHasPageName = false
	)
	// Check the page variable in the URI.
	if len(r.Router.RegNames) > 0 {
		for _, name := range r.Router.RegNames {
			if name == gpage.DefaultPageName {
				uriHasPageName = true
				break
			}
		}
		if uriHasPageName {
			if match, err := gregex.MatchString(r.Router.RegRule, url.Path); err == nil && len(match) > 0 {
				if len(match) > len(r.Router.RegNames) {
					urlTemplate = r.Router.Uri
					for i, name := range r.Router.RegNames {
						rule := fmt.Sprintf(`[:\*]%s|\{%s\}`, name, name)
						if name == gpage.DefaultPageName {
							urlTemplate, err = gregex.ReplaceString(rule, gpage.DefaultPagePlaceHolder, urlTemplate)
						} else {
							urlTemplate, err = gregex.ReplaceString(rule, match[i+1], urlTemplate)
						}
						if err != nil {
							panic(err)
						}
					}
				}
			} else {
				panic(err)
			}
		}
	}
	// Check the page variable in the query string.
	if !uriHasPageName {
		values := url.Query()
		values.Set(gpage.DefaultPageName, gpage.DefaultPagePlaceHolder)
		url.RawQuery = values.Encode()
		// Replace the encoded "{.page}" to original "{.page}".
		url.RawQuery = gstr.Replace(url.RawQuery, "%7B.page%7D", "{.page}")
	}
	if url.RawQuery != "" {
		urlTemplate += "?" + url.RawQuery
	}

	return gpage.New(totalSize, pageSize, r.Get(gpage.DefaultPageName).Int(), urlTemplate)
}
