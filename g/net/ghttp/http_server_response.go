package ghttp

import (
    "net/http"
    "gitee.com/johng/gf/g/encoding/gjson"
    "sync"
)

// 服务端请求返回对象
type ServerResponse struct {
    http.ResponseWriter
    bufmu  sync.RWMutex // 缓冲区互斥锁
    buffer []byte       // 每个请求的返回数据缓冲区
}

// 返回的固定JSON数据结构
type ResponseJson struct {
    Result  int         `json:"result"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 返回信息(byte)
func (r *ServerResponse) Write(content []byte) {
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    r.buffer = append(r.buffer, content...)
}

// 返回信息(string)
func (r *ServerResponse) WriteString(content string) {
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    r.buffer = append(r.buffer, content...)
}

// 返回固定格式的json
func (r *ServerResponse) WriteJson(result int, message string, data interface{}) error {
    r.Header().Set("Content-Type", "application/json")
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    if jsonstr, err := gjson.Encode(ResponseJson{ result, message, data }); err != nil {
        return err
    } else {
        r.buffer = append(r.buffer, jsonstr...)
    }
    return nil
}

// 返回内容编码
func (r *ServerResponse) WriteHeaderEncoding(encoding string) {
    r.Header().Set("Content-Type", "text/plain; charset=" + encoding)
}

// 获取缓冲区数据
func (r *ServerResponse) Buffer() []byte {
    r.bufmu.RLock()
    defer r.bufmu.RUnlock()
    return r.buffer
}

// 输出缓冲区数据到客户端
func (r *ServerResponse) Output() {
    r.bufmu.RLock()
    defer r.bufmu.RUnlock()
    r.ResponseWriter.Write(r.buffer)
}
