package gluster

import (
    "time"
    "g/util/gtime"
    "strings"
    "sync"
    "g/database/gdb"
    "g/net/ghttp"
    "g/core/types/gmap"
    "g/os/glog"
    "fmt"
    "strconv"
    "g/os/gcache"
)

// 将Service转为可json化的数据结构
func (n *Node) serviceMapToServiceStructMap() *map[string]ServiceStruct {
    stm := make(map[string]ServiceStruct)
    for k, v := range *n.Service.Clone() {
        s     := v.(Service)
        stm[k] = *n.serviceToServiceStruct(&s)
    }
    return &stm
}

// 将Service对象转换为ServiceStruct对象，以便json化
func (n *Node) serviceToServiceStruct(s *Service) *ServiceStruct {
    var st ServiceStruct
    st.Name = s.Name
    st.Type = s.Type
    st.Node = make(map[string]interface{})
    for k, v := range *s.Node.Clone() {
        m := v.(*gmap.StringInterfaceMap)
        st.Node[k] = *m.Clone()
    }
    return &st
}

// 将ServiceStruct对象转换为Service对象
func (n *Node) serviceSructToService(st *ServiceStruct) *Service {
    var s Service
    s.Name = st.Name
    s.Type = st.Type
    s.Node = gmap.NewStringInterfaceMap()
    for k, v := range st.Node {
        m := gmap.NewStringInterfaceMap()
        m.BatchSet(v.(map[string]interface{}))
        s.Node.Set(k, m)
    }
    return &s
}

// 服务健康检查回调函数
func (n *Node) serviceHealthCheckHandler() {
    start := gtime.Millisecond()
    for {
        if n.getRaftRole() == gROLE_RAFT_LEADER && gtime.Millisecond() > start {
            for _, v := range n.Service.Values() {
                service := v.(Service)
                go func(s *Service) {
                    n.checkServiceHealth(s)

                }(&service)
            }
            start = gtime.Millisecond() + gSERVICE_HEALTH_CHECK_INTERVAL
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// 服务健康检测
// 如果新增检测类型，需要更新该方法
func (n *Node) checkServiceHealth(service *Service) {
    var wg sync.WaitGroup
    // 用以标识本分组的Service是否有更新
    updated := false
    for k, v := range *service.Node.Clone() {
        wg.Add(1)
        go func(name string, m *gmap.StringInterfaceMap, u *bool) {
            cachekey  := "gluster_service_" + service.Name + "_" + name + "_check"
            needcheck := gcache.Get(cachekey)
            if needcheck != nil {
                return
            }
            //glog.Printf("start checking node: %s, name: %s, \n", name, service.Name)
            ostatus := m.Get("status")
            switch strings.ToLower(service.Type) {
                case "mysql": fallthrough
                case "pgsql": n.dbHealthCheck(service.Type, m)
                case "web":   n.webHealthCheck(m)
            }
            nstatus := m.Get("status")
            if ostatus != nstatus {
                (*u) = true
                glog.Printf("service updated, node: %s, from %v to %v, name: %s, \n", name, ostatus, nstatus, service.Name)
            }
            interval := m.Get("interval")
            timeout  := int64(gSERVICE_HEALTH_CHECK_INTERVAL)
            if interval != nil {
                timeout, _ = strconv.ParseInt(interval.(string), 10, 64)
            }
            gcache.Set(cachekey, 1, timeout)
            wg.Done()
        }(k, v.(*gmap.StringInterfaceMap), &updated)
    }
    wg.Wait()
    // 从Service对象为基础，新创建一个ServiceStruct，更新到API接口变量中，以便提高接口查询效率
    // 以空间换时间的方式，在ServiceForApi中已有的变量会被新变量替换，但旧变量不会马上消失，而是转交给GC处理
    if updated {
        n.ServiceForApi.Set(service.Name, *n.serviceToServiceStruct(service))
        n.setLastServiceLogId(gtime.Microsecond())
    }
}

// MySQL/PostgreSQL数据库健康检查
// 使用并发方式并行测试同一个配置中的数据库链接
func (n *Node) dbHealthCheck(stype string, item *gmap.StringInterfaceMap) {
    host := item.Get("host")
    port := item.Get("port")
    user := item.Get("user")
    pass := item.Get("pass")
    name := item.Get("database")
    if host == nil || port == nil || user == nil || pass == nil || name == nil {
        return
    }
    dbcfg   := gdb.ConfigNode{
        Host    : host.(string),
        Port    : port.(string),
        User    : user.(string),
        Pass    : pass.(string),
        Name    : name.(string),
        Type    : stype,
    }
    db, err := gdb.NewByConfigNode(dbcfg)
    if err != nil || db == nil {
        item.Set("status", 0)
    } else {
        if db.PingMaster() != nil {
            item.Set("status", 0)
        } else {
            item.Set("status", 1)
        }
        db.Close()
    }
}

// WEB健康检测
func (n *Node) webHealthCheck(item *gmap.StringInterfaceMap) {
    url   := ""
    check := item.Get("check")
    if check == nil {
        host := item.Get("host")
        port := item.Get("port")
        if host != nil && port != nil {
            url = fmt.Sprintf("http://%s:%s", host, port)
        }
    } else {
        url = check.(string)
    }
    if url == "" {
        return
    }
    r := ghttp.Get(url)
    if r == nil || r.StatusCode != 200 {
        item.Set("status", 0)
    } else {
        item.Set("status", 1)
    }

    if r != nil {
        r.Close()
    }
}