package graft

import (
    "time"
    "g/util/gtime"
    "strings"
    "sync"
    "g/database/gdb"
)

func (n *Node) serviceHealthCheckHandler() {
    var wg sync.WaitGroup
    start := gtime.Millisecond()
    for {
        if n.getRole() == gROLE_LEADER && gtime.Millisecond() >= start {
            start += gSERVICE_HEALTH_CHECK_INTERVAL
            for k, v := range *n.Service.Clone() {
                wg.Add(1)
                go func(service *Service) {
                    switch strings.ToLower(service.Type) {
                        case "mysql": n.mysqlHealthCheck(k, service)
                        case "web":   n.webHealthCheck(k, service)
                    }
                    wg.Done()
                }(&(v.(Service)))
            }
            wg.Wait()
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func (n *Node) updateService(name string, service *Service) {
    for _, m := range service.List {
        dbcfg := gdb.ConfigNode{
            Host    : m["host"].(string),
            Port    : m["port"].(int),
            User    : m["user"].(string),
            Pass    : m["pass"].(string),
            Name    : m["database"].(string),
            Type    : "mysql",
            Role    : "master",
        }
        db, err := gdb.NewByConfigNode(dbcfg)
        if err != nil {
            m["status"] = 0
        } else {

        }
    }
}

func (n *Node) mysqlHealthCheck(name string, service *Service) {
    for _, m := range service.List {

    }
}

func (n *Node) webHealthCheck(name string, service *Service) {

}