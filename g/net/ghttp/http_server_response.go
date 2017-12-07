package ghttp

import (
    "net/http"
    "gitee.com/johng/gf/g/encoding/gjson"
)

// 服务端请求返回对象
type ServerResponse struct {
    http.ResponseWriter
    server *Server      // 所属Server对象
}

// 返回的固定JSON数据结构
type ResponseJson struct {
    Result  int         `json:"result"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 返回信息
func (r *ServerResponse) Write(content []byte) {
    if r.Header().Get("Content-Type") == "" {
        r.Header().Set("Content-Type", "text/plain; charset=utf-8")
    }
    r.ResponseWriter.Write(content)
}

// 返回固定格式的json
func (r *ServerResponse) ResponseJson(result int, message string, data interface{}) {
    if r.Header().Get("Content-Type") == "" {
        r.Header().Set("Content-Type", "application/json")
    }
    r.Write([]byte(gjson.Encode(ResponseJson{ result, message, data })))
}

