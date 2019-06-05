package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type Session struct {
	Token string `json:"token,omitempty"`
	Ttl   int    `json:"ttl,omitempty"` //token 存活时间，单位：秒
}

type Message struct {
	Code  int    `json:"code"`
	Body  string `json:"body,omitempty"`
	Error string `json:"error,omitempty"`
}

type Paginate struct {
	CurrentPage    string `json:"current_page,omitempty"`
	PrePage        string `json:"pre_page,omitempty"`
	NextPage       string `json:"next_page,omitempty"`
	LastPage       string `json:"last_page,omitempty"`
	PageSize       string `json:"page_size,omitempty"`
	TotalCount     string `json:"total_count,omitempty"`
	TotalPage      string `json:"total_page,omitempty"`
	PrePageUrl     string `json:"pre_page_url,omitempty"`
	NextPageUrl    string `json:"next_page_url,omitempty"`
	CurrentPageUrl string `json:"current_page_url,omitempty"`
}

type ResponseJson struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	ExtData  interface{} `json:"ext_data,omitempty"`
	Paginate interface{} `json:"paginate,omitempty"`
	Message  Message     `json:"message,omitempty"`
}

// 错误json
func ErrorJson(r *ghttp.Request, message Message) {
	responseJson := &ResponseJson{
		Success: false,
		Message: message,
	}
	//responseJson := struct {
	//	Code string
	//}{
	//	Code:"ok",
	//}
	fmt.Println(responseJson)
	_ = r.Response.WriteJson(responseJson)
	//glog.Error(error.Error())
	r.Exit()
}

// 成功但没有数据的json
func SuccessJson(r *ghttp.Request, message Message) {
	responseJson := &ResponseJson{
		Success: true,
		Message: message,
	}
	_ = r.Response.WriteJson(responseJson)
	//glog.Error(error.Error())
	r.Exit()
}

// 成功但带有数据的json
func DataJson(r *ghttp.Request, success bool, message Message, data interface{}, extData interface{}) {
	responseJson := &ResponseJson{
		Success: true,
		Data:    data,
		ExtData: extData,
		Message: message,
	}
	_ = r.Response.WriteJson(responseJson)
	//glog.Error(error.Error())
	r.Exit()
}

// 成功但带有数据和Session的json
func SessionJson(r *ghttp.Request, success bool, message Message, data interface{}, session Session) {
	responseJson := &ResponseJson{
		Success: success,
		Data:    data,
		ExtData: session,
		Message: message,
	}
	_ = r.Response.WriteJson(responseJson)
	//glog.Error(error.Error())
	r.Exit()
}

// 成功带有数据和分页信息的json
func PaginateJson(r *ghttp.Request, success bool, message Message, paginate Paginate, data interface{}, extData interface{}) {
	responseJson := &ResponseJson{
		Success:  success,
		Data:     data,
		ExtData:  extData,
		Paginate: paginate,
		Message:  message,
	}
	_ = r.Response.WriteJson(responseJson)
	//glog.Error(error.Error())
	r.Exit()
}


func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.BindHandler("/", func(r *ghttp.Request) {
		DataJson(r, true, Message{3, "测试", ""}, nil, nil)
	})
	s.SetPort(8199)
	s.Run()
}
