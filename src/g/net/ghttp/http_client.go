package ghttp

import (
    "net/http"
    "log"
    "strings"
    "io/ioutil"
    "time"
)

// 设置请求过期时间
func (c *Client) SetTimeOut(t time.Duration)  {
    c.Timeout = t
}

// post请求提交数据
func (c *Client) Post(url, data string) string {
    resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
    if err != nil {
        log.Println(err)
        return ""
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
        resp.Body.Close()
        return ""
    }
    resp.Body.Close()
    return string(body)
}

// get请求
func (c *Client) Get(url string) string {
    resp, err := http.Get(url)
    if err != nil {
        log.Println(err)
        return ""
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Println(err)
        resp.Body.Close()
        return ""
    }
    resp.Body.Close()
    return string(body)
}

// 请求并返回response对象
func (c *Client) Request(method, url, data string) *Response {
    client   := &http.Client{}
    req, err := http.NewRequest(method, url, strings.NewReader(data))
    if err != nil {
        log.Println(err)
        return nil
    }
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return nil
    }
    r := &Response{}
    r.Response = *resp
    return r
}

