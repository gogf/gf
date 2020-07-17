package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
)

var (
	tableName = "orders"
	dao       = g.DB().Table(tableName).Safe()
)

type OrderServiceEntity struct {
	GoodsPrice float64 `json:"goods_price" gvalid:"required"`
	PayTo      int8    `json:"payTo" gvalid:"required"`
	PayStatus  int8    `json:"payStatus" `
	CreateTime string  `json:"createTime" `
	AppId      string  `json:"appId" gvalid:"required"`
	PayUser    string  `json:"pay_user" gvalid:"required"`
	QrUrl      string  `json:"qr_url" `
}

type Create struct {
	Id         int64   `json:"id" gconv:"id"`
	GoodsPrice float64 `json:"goodsPrice" gconv:"goods_price"`
	PayTo      int8    `json:"payTo" gconv:"pay_to"`
	PayStatus  int8    `json:"payStatus" gconv:"pay_status"`
	CreateTime string  `json:"createTime" gconv:"create_time"`
	UserId     int     `json:"user_id" `
	PayUser    string  `json:"pay_user" `
	QrUrl      string  `json:"qr_url" `
}

func main() {
	g.DB().SetDebug(true)
	userInfo := Create{
		Id: 3,
	}
	orderService := OrderServiceEntity{
		GoodsPrice: 0.1,
		PayTo:      1,
	}
	size, err := dao.Where("user_id", userInfo.Id).
		And("goods_price", float64(100.10)).
		And("pay_status", 0).
		And("pay_to", orderService.PayTo).Count()
	fmt.Println(err)
	fmt.Println(size)
}
