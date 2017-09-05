package gluster

import (
    "g/net/ghttp"
    "fmt"
    "errors"
    "g/util/grand"
    "strconv"
    "g/os/gcache"
)

// 用于负载均衡计算的结构体
type PriorityNode struct {
    name     string
    priority int
}

// 负载均衡查询
func (this *NodeApiBalance) GET(r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    name := r.GetRequestString("name")
    if name == "" {
        w.ResponseJson(0, "incomplete input: name is required", nil)
    } else {
        // 高并发下的缓存处理，缓存时间为1秒
        k := "gluster_service_balance_name_" + name
        r := gcache.Get(k)
        if r == nil {
            r, err := this.getAliveServiceByPriority(name)
            if err != nil {
                w.ResponseJson(0, err.Error(), nil)
            } else {
                gcache.Set(k, r, 1)
                w.ResponseJson(0, "ok", r)
            }
        } else {
            w.ResponseJson(0, "ok", r)
        }
    }
}

// 查询存货的service, 并根据priority计算负载均衡，取出一条返回
func (this *NodeApiBalance) getAliveServiceByPriority(name string ) (interface{}, error) {
    if !this.node.ServiceForApi.Contains(name) {
        return nil, errors.New(fmt.Sprintf("no service named '%s' found", name))
    }
    st   := this.node.ServiceForApi.Get(name).(ServiceStruct)
    list := make([]PriorityNode, 0)
    for k, v := range st.Node {
        m := v.(map[string]interface{})
        status, ok := m["status"]
        if !ok || status.(int) == 0 {
            continue
        }
        priority, ok := m["priority"]
        if !ok {
            continue
        }
        r, err := strconv.Atoi(priority.(string))
        if err == nil {
            list = append(list, PriorityNode{k, r})
        }
    }
    if len(list) < 1 {
        return nil, errors.New("no nodes are alive")
    }
    nodename := this.getServiceByPriority(list)
    if nodename == "" {
        return nil, errors.New("get node by balance failed, please check the data structure of the service")
    }
    return st.Node[nodename], nil
}

// 根据priority计算负载均衡
func (this *NodeApiBalance) getServiceByPriority (list []PriorityNode) string {
    if len(list) < 2 {
        return list[0].name
    }
    var total int
    for i := 0; i < len(list); i++ {
        total += list[i].priority * 100
    }
    r   := grand.Rand(0, total)
    min := 0
    max := 0
    for i := 0; i < len(list); i++ {
        max = min + list[i].priority * 100
        //fmt.Printf("r: %d, min: %d, max: %d\n", r, min, max)
        if r >= min && r < max {
            return list[i].name
        } else {
            min = max
        }
    }
    return ""
}
