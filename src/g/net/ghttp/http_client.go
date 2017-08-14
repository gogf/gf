package ghttp

import (
    "net/http"
    "log"
    "strings"
    "time"
)

// 设置请求过期时间
func (c *Client) SetTimeOut(t time.Duration)  {
    c.Timeout = t
}

// GET请求
func (c *Client) Get(url string) *ClientResponse {
    return c.Request("GET", url, "")
}

// PUT请求
func (c *Client) Put(url, data string) *ClientResponse {
    return c.Request("PUT", url, data)
}

// POST请求提交数据
func (c *Client) Post(url, data string) *ClientResponse {
    resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
    if err != nil {
        log.Println(err)
        return nil
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r
}

// DELETE请求
func (c *Client) Delete(url, data string) *ClientResponse {
    return c.Request("DELETE", url, data)
}

// 请求并返回response对象
func (c *Client) Request(method, url, data string) *ClientResponse {
    client   := &http.Client{}
    req, err := http.NewRequest(strings.ToUpper(method), url, strings.NewReader(data))
    if err != nil {
        log.Println(err)
        return nil
    }
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return nil
    }
    r := &ClientResponse{}
    r.Response = *resp
    return r
}

