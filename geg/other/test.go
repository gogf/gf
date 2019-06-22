package main

import (
	"fmt"

	"github.com/gogf/gf/g/encoding/gjson"
	"github.com/gogf/gf/g/util/gconv"
)

//type TemplateMessage struct {
//	Touser      string      `json:"touser,omitempty"`
//	TemplateId  string      `json:"template_id,omitempty"`
//	Miniprogram interface{} `json:"miniprograme,omitempty"`
//	Data        interface{} `json:"data,omitempty"`
//}

type TemplateMessage struct {
	Touser      string      `json:"touser,omitempty"`
	TemplateId  string      `json:"template_id,omitempty"`
	Miniprogram *gjson.Json `json:"miniprograme,omitempty"`
	Data        *gjson.Json `json:"data,omitempty"`
}

// 封装模版消息
func getTemplateMessage(message string) string {
	templateId := "22222222"
	miniprogram := `{"appid":"111111111","pagepath":"pages\/index?ald_media_id=20962&ald_link_key=bd660b4962a599f2"}`
	data := `{"first":{"value":"送您一个随机红包，点击领取¥0.3-¥10","color":"#FF0000"},"keyword1":{"value":"¥0.3-¥10"},"keyword2":{"value":"2019年06月19日 11:45"},"keyword3":{"value":"微信零钱"},"keyword4":{"value":"点击此消息即可提现","color":"#FF0000"}}`
	miniprogramJson := gjson.New(miniprogram)
	dataJson := gjson.New(data)
	//glog.Infof(miniprogramJson.ToJsonString())
	//glog.Info(dataJson.ToJsonString())

	templateMessage := TemplateMessage{
		Touser:      message,
		TemplateId:  templateId,
		Miniprogram: miniprogramJson,
		Data:        dataJson,
	}

	//glog.Debug(dataJson.ToJsonString())
	//json, _ := gjson.New(templateMessage).ToJsonString()
	//return json
	//glog.Info(templateMessage.Miniprogram.ToJsonString())
	return gconv.String(templateMessage)
}

func main() {
	fmt.Println(getTemplateMessage("test"))

}
