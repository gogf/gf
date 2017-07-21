package gip

import (
    "net"
    "strconv"
    "strings"
)

// 获取本地局域网ip列表
func IntranetIP() (ips []string, err error) {
    ips        = make([]string, 0)
    ifaces, e := net.Interfaces()
    if e != nil {
        return ips, e
    }

    for _, iface := range ifaces {
        if iface.Flags&net.FlagUp == 0 {
            continue // interface down
        }

        if iface.Flags & net.FlagLoopback != 0 {
            continue // loopback interface
        }

        // ignore warden bridge
        if strings.HasPrefix(iface.Name, "w-") {
            continue
        }

        addrs, e := iface.Addrs()
        if e != nil {
            return ips, e
        }

        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
                case *net.IPNet:
                    ip = v.IP
                case *net.IPAddr:
                    ip = v.IP
            }

            if ip == nil || ip.IsLoopback() {
                continue
            }

            ip = ip.To4()
            if ip == nil {
                continue // not an ipv4 address
            }

            ipStr := ip.String()
            if IsIntranet(ipStr) {
                ips = append(ips, ipStr)
            }
        }
    }
    return ips, nil
}

// 判断所给ip是否为局域网ip
// A类 10.0.0.0--10.255.255.255
// B类 172.16.0.0--172.31.255.255
// C类 192.168.0.0--192.168.255.255
func IsIntranet(ipStr string) bool {
    // ip协议保留的局域网ip
    if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
        return true
    }
    if strings.HasPrefix(ipStr, "172.") {
        // 172.16.0.0 - 172.31.255.255
        arr := strings.Split(ipStr, ".")
        if len(arr) != 4 {
            return false
        }

        second, err := strconv.ParseInt(arr[1], 10, 64)
        if err != nil {
            return false
        }

        if second >= 16 && second <= 31 {
            return true
        }
    }

    return false
}