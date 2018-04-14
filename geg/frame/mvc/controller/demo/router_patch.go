package demo

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/router-patch", RouterPatch)
    ghttp.GetServer().Router.SetPatchRule(`\/list\?page=(\d+)&*`, "/list/page/$1?")
}

func RouterPatch(r *ghttp.Request) {
    r.Response.Write(`<a href="/list?page=2&ajax=1">page2</a>`)
}