// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// 获得当前请求，指定类型的的钩子函数列表
func (r *Request) getHookHandlers(hook string) []*handlerParsedItem {
	if !r.hasHookHandler {
		return nil
	}
	parsedItems := make([]*handlerParsedItem, 0, 4)
	for _, v := range r.handlers {
		if v.handler.hookName != hook {
			continue
		}
		item := v
		parsedItems = append(parsedItems, item)
	}
	return parsedItems
}
