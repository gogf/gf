package gluster

import (
    "time"
    "g/util/gtime"
    "strings"
    "sync"
    "g/database/gdb"
    "g/net/ghttp"
    "g/core/types/gmap"
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
    st.List = make([]map[string]interface{}, len(s.List))
    for i, v2 := range s.List {
        st.List[i] = *v2.Clone()
    }
    return &st
}

// 将ServiceStruct对象转换为Service对象
func (n *Node) serviceSructToService(st *ServiceStruct) *Service {
    var s Service
    s.Name = st.Name
    s.Type = st.Type
    s.List = make([]*gmap.StringInterfaceMap, len(st.List))
    for i, v2 := range st.List {
        m := gmap.NewStringInterfaceMap()
        m.BatchSet(v2)
        s.List[i] = m
    }
    return &s
}

// 服务健康检查回调函数
func (n *Node) serviceHealthCheckHandler() {
    var wg sync.WaitGroup
    start := gtime.Millisecond()
    for {
        if n.getRaftRole() == gROLE_RAFT_LEADER && gtime.Millisecond() > start {
            for _, v := range n.Service.Values() {
                service := v.(Service)
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
    var wg sync.WaitGroup
    // 用以标识Service是否有更新
    updated := false
    for k, v := range service.List {
        wg.Add(1)
        go func(i int, m *gmap.StringInterfaceMap, u *bool) {
            ostatus := m.Get("status")
            switch strings.ToLower(service.Type) {
                case "mysql": fallthrough
                case "pgsql": n.dbHealthCheck(service.Type, m)
                case "web":   n.webHealthCheck(m)
            }
            if ostatus != m.Get("status") {
                (*u) = true
            }
            wg.Done()
        }(k, v, &updated)
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
    dbcfg   := gdb.ConfigNode{
        Host    : item.Get("host").(string),
        Port    : item.Get("port").(string),
        User    : item.Get("user").(string),
        Pass    : item.Get("pass").(string),
        Name    : item.Get("database").(string),
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
    url := item.Get("check")
    if url == nil {
        url = item.Get("url")
    }
    r := ghttp.Get(url.(string))
    if r == nil || r.StatusCode != 200 {
        item.Set("status", 0)
    } else {
        item.Set("status", 1)
    }

    if r != nil {
        r.Close()
    }
}