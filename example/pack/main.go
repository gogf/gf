// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package main

import (
	_ "github.com/wangyougui/gf/example/pack/packed"

	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/i18n/gi18n"
	"github.com/wangyougui/gf/v2/net/ghttp"
	"github.com/wangyougui/gf/v2/os/gres"
)

func main() {
	gres.Dump()

	s := g.Server()
	s.SetPort(8199)
	s.SetServerRoot("resource/public")
	s.BindHandler("/i18n", func(r *ghttp.Request) {
		var (
			lang    = r.Get("lang", "zh-CN").String()
			ctx     = gi18n.WithLanguage(r.Context(), lang)
			content string
		)
		content = g.I18n().T(ctx, `{#hello} {#world}!`)
		r.Response.Write(content)
	})
	s.Run()
}
