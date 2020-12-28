package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	var (
		orderId     = 865271654
		orderAmount = 99.8
	)
	fmt.Println(g.I18n().Tfl(`en`, `{#OrderPaid}`, orderId, orderAmount))
	fmt.Println(g.I18n().Tfl(`zh-CN`, `{#OrderPaid}`, orderId, orderAmount))
}
