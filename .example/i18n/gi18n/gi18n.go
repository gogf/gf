package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
)

func main() {
	var (
		orderId     = 865271654
		orderAmount = 99.8
	)
	fmt.Println(g.I18n().Tf(
		gi18n.WithLanguage(context.TODO(), `en`),
		`{#OrderPaid}`, orderId, orderAmount,
	))
	fmt.Println(g.I18n().Tf(
		gi18n.WithLanguage(context.TODO(), `zh-CN`),
		`{#OrderPaid}`, orderId, orderAmount,
	))
}
