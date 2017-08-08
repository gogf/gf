package ghttp

import (
    "g/encoding/gjson"
    "io/ioutil"
    "log"
)

type ResponseJson struct {
    Result  int         `json:"result"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 返回固定格式的json
func (r *Response) ResponseJson(result int, message string, data interface{}) {
    r.writer.Header().Set("Content-type", "application/json")
    r.writer.Write([]byte(*gjson.Encode(ResponseJson{ result, message, data })))
}

// 获取返回的数据
func (r *Response) ReadAll() string {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println(err)
        r.Body.Close()
        return ""
    }
    r.Body.Close()
    return string(body)
}