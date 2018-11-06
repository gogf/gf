package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    s := `
172.20.1.198 - - [2018-11-06T16:26:09+08:00] "POST /passport HTTP/1.1" "OK" 1 200 0.000 0.035 0.035 0.035 448 "-" "-" "-" "{\x22jsonrpc\x22:\x222.0\x22,\x22method\x22:\x22getSessionInfo\x22,\x22params\x22:[\x2262703819__6augmxzV9f5c7o4MEimnMqPhoyKWPi8pXjs2VIj3T43vBfuGZOJ9DxrbRNsFB0ew\x22,true,{\x22platform\x22:\x22web-ph\x22}],\x22id\x22:1}" http unix:/var/run/php/php5.6-fpm.sock med3-svr [med3-svr-65494945bf-jppth] "" "" "" [med3-svr-65494945bf-jppth]
`

    fmt.Println(gtime.ParseTimeFromContent(s))
}