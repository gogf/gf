package main

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

//
type BaseData struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	BeginTime int64  `json:"begin_time"`
	EndTime   int64  `json:"end_time"`
}

// 奖励数据
type RewardOneData struct {
	Type  string `json:"type"`
	Id    uint32 `json:"id"`
	Count uint32 `json:"count"`
}

// 邮件
type MailOneData struct {
	Base       BaseData        `json:"base"`
	MailId     uint32          `json:"mail_id"`
	PicId      uint32          `json:"pic_id"`
	PlayerList g.ArrayStr      `json:"player_list"`
	Reward     []RewardOneData `json:"reward"`
}

func main() {
	test()
}

func test() {
	// 使用下面这行就会panic 原因是 最后的reward字段列表为空
	jsonStr := "{\"base\":{\"title\":\"testTitle\",\"content\":\"testContent\",\"begin_time\":1574763804,\"end_time\":1574767404},\"mail_id\":1,\"pic_id\":1,\"player_list\":[],\"reward\":[]}"
	//jsonStr := "{\"base\":{\"title\":\"testTitle\",\"content\":\"testContent\",\"begin_time\":1574763804,\"end_time\":1574767404},\"mail_id\":1,\"pic_id\":1,\"player_list\":[],\"reward\":[{\"type\":\"diamond\",\"id\":0,\"count\":100}]}"
	decodeData, _ := gjson.Decode(jsonStr)

	decodeMailData := new(MailOneData)
	gconv.Struct(decodeData, decodeMailData)

	//
	g.Dump(decodeMailData)
}
