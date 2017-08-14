package gluster

import (
    "time"
    "g/util/gtime"
    "strings"
    "sync"
    "g/database/gdb"
    "g/net/ghttp"
)

// 服务健康检查回调函数
func (n *Node) serviceHealthCheckHandler() {
    var wg sync.WaitGroup
    start := gtime.Millisecond()
    for {
        if n.getRole() == gROLE_LEADER && gtime.Millisecond() > start {
            for _, v := range n.Service.Values() {
                service := v.(Service)
                //log.Println("health check:", service.Name)
                go func(s *Service) {
                    wg.Add(1)
                    n.checkServiceHealth(s)
                    wg.Done()
                }(&service)
            }
            wg.Wait()
            start = gtime.Millisecond() + gSERVICE_HEALTH_CHECK_INTERVAL
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// 服务健康检测
// 如果新增检测类型，需要更新该方法
func (n *Node) checkServiceHealth(service *Service) {
    var wg    sync.WaitGroup
    var mutex sync.Mutex
    newList := make([]map[string]interface{}, len(service.List))
    for k, v := range service.List {
        wg.Add(1)
        go func(i int, m map[string]interface{}) {
            item := make(map[string]interface{})
            switch strings.ToLower(service.Type) {
                case "mysql": fallthrough
                case "pgsql": item = n.dbHealthCheck(service.Type, m)
                case "web":   item = n.webHealthCheck(m)
            }
            mutex.Lock()
            newList[i] = item
            mutex.Unlock()
            wg.Done()
        }(k, v)
    }
    wg.Wait()
    // 替换服务列表为最新状态列表
    service.List = newList
    // 保存检测结果到Service成员变量
    n.Service.Set(service.Name, *service)
}

// MySQL/PostgreSQL数据库健康检查
// 使用并发方式并行测试同一个配置中的数据库链接
func (n *Node) dbHealthCheck(stype string, item map[string]interface{}) map[string]interface{} {
    dbcfg := gdb.ConfigNode{
        Host    : item["host"].(string),
        Port    : item["port"].(string),
        User    : item["user"].(string),
        Pass    : item["pass"].(string),
        Name    : item["database"].(string),
        Type    : stype,
    }
    db, err := gdb.NewByConfigNode(dbcfg)
    if err != nil || db == nil {
        item["status"] = 0
    } else {
        if db.PingMaster() != nil {
            item["status"] = 0
        } else {
            item["status"] = 1
        }
        db.Close()
    }
    return item
}

// WEB健康检测
func (n *Node) webHealthCheck(item map[string]interface{}) map[string]interface{} {
    url, ok := item["check"]
    if !ok {
        url = item["url"]
    }
    c := ghttp.NewClient()
    r := c.Get(url.(string))
    if r == nil || r.StatusCode != 200 {
        item["status"] = 0
    } else {
        item["status"] = 1
    }
    return item
}