package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		rs :="GF(Go Frame)是一款模块化、松耦合、生产级的Go应用开发框架。提供了常用的核心开发组件，如：缓存、日志、文件、时间、队列、数组、集合、字符串、定时器、命令行、文件锁、内存锁、对象池、连接池、数据校验、数据编码、文件监控、定时任务、数据库ORM、TCP/UDP组件、进程管理/通信、 并发安全容器等等。并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、路由注册、配置管理、模板引擎等等，支持热重启、热更新、多域名、多端口、多服务、HTTPS、Rewrite等特性。"
		//此行压测会提示map并发错误   webbench -c 8000 -t 60  http://IP  局域网两台机器测试
		r.Response.WriteTplContent(rs, g.Map{
			"Contentb": 1,
		})
	})
	s.SetPort(8199)
	s.Run()
}