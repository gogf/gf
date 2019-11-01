package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gipv4"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"math"
	"net"
)

func main() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		fmt.Println(net.ParseMAC(netInterface.HardwareAddr.String()))
	}

	return
	ip, _ := gipv4.IntranetIP()

	g.Dump(math.MaxInt64)
	g.Dump(gipv4.Ip2long(ip), gbase64.EncodeString(ip))
	g.Dump(gproc.Pid())
	g.Dump(gtime.Nanosecond())
}
