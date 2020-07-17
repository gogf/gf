package types

import "github.com/jin502437344/gf/frame/g"

type NodeInfo struct {
	Cpu      float32 // CPU百分比(%)
	Host     string  // 主机名称
	Ip       g.Map   // IP地址信息(可能多个)
	MemUsed  int     // 内存使用(byte)
	MemTotal int     // 内存总量(byte)
	Time     int     // 上报时间(时间戳)
}
