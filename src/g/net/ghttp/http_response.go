package ghttp

import (
    "g/encoding/gjson"
    "io/ioutil"
    "g/os/glog"
)

type ResponseJson struct {
    Result  int         `json:"result"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 返回固定格式的json
func (r *ServerResponse) ResponseJson(result int, message string, data interface{}) {
    r.Header().Set("Content-type", "application/json")
    r.Write([]byte(*gjson.Encode(ResponseJson{ result, message, data })))
}

// 获取返回的数据
func (r *ClientResponse) ReadAll() string {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        glog.Println(err)
        r.Body.Close()
        return ""
    }
    r.Body.Close()
    return string(body)
}