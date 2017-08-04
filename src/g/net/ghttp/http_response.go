package ghttp

import (
    "g/encoding/gjson"
    "io"
)

type ResponseJson struct {
    Result  int         `json:"result"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func (r *Response) ResponseJson(result int, message string, data interface{}) {
    io.WriteString(r, *gjson.Encode(ResponseJson{ result, message, data }))
}