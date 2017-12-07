package ghttp

import (
    "io/ioutil"
)

// 获取返回的数据
func (r *ClientResponse) ReadAll() []byte {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return nil
    }
    return body
}

// 关闭返回的HTTP链接
func (r *ClientResponse) Close()  {
    r.Response.Close = true
    r.Body.Close()
}