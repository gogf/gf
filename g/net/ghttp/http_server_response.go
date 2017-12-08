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

// 返回信息(byte)
func (r *ServerResponse) Write(content []byte) {
    r.ResponseWriter.Write(content)
}

// 返回信息(string)
func (r *ServerResponse) WriteString(content string) {
    r.Write([]byte(content))
}

// 返回固定格式的json
func (r *ServerResponse) WriteJson(result int, message string, data interface{}) {
    if r.Header().Get("Content-Type") == "" {
        r.Header().Set("Content-Type", "application/json")
    }
    r.Write([]byte(gjson.Encode(ResponseJson{ result, message, data })))
}

// 返回内容编码
func (r *ServerResponse) WriteHeaderEncoding(encoding string) {
    r.Header().Set("Content-Type", "text/plain; charset=" + encoding)
}

