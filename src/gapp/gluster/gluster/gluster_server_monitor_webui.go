package gluster

import (
    "g/net/ghttp"
)


// service 查询
func (this *MonitorWebUI) GET(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.ResponseJson(1, "ok", nil)
}



