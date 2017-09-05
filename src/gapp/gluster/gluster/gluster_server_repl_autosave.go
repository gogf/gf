// 数据同步需要注意的是：
// leader只有在通知完所有follower更新完数据之后，自身才会进行数据更新
// 因此leader
package gluster

import (
    "g/encoding/gjson"
    "time"
    "g/os/gfile"
    "g/util/gtime"
    "g/os/glog"
    "g/encoding/gcompress"
)

// 日志自动保存处理
func (n *Node) logAutoSavingHandler() {
    t := gtime.Millisecond()
    for {
        // 当日志列表的最新ID与保存的ID不相等，或者超过超时时间
        if n.getLastLogId() != n.getLastSavedLogId() || gtime.Millisecond() - t > gLOG_REPL_AUTOSAVE_INTERVAL {
            //glog.Println("saving data to file")
            n.saveDataToFile()
            t = gtime.Millisecond()
        } else {
            time.Sleep(100 * time.Millisecond)
        }
    }
}

// 保存数据到磁盘
func (n *Node) saveDataToFile() {
    data := SaveInfo {
        LastLogId           : n.getLastLogId(),
        LogCount            : n.getLogCount(),
        LogList             : make([]LogEntry, 0),
        LastServiceLogId    : n.getLastServiceLogId(),
        Service             : *n.serviceMapToServiceStructMap(),
        Peers               : *n.Peers.Clone(),
        DataMap             : *n.DataMap.Clone(),
    }
    for _, v := range n.LogList.BackAll() {
        data.LogList = append(data.LogList, v.(LogEntry))
    }
    content := []byte(gjson.Encode(&data))
    if gCOMPRESS_SAVING {
        content = gcompress.Zlib(content)
    }
    err := gfile.PutBinContents(n.getDataFilePath(), content)
    if err != nil {
        glog.Error("saving data error:", err)
    } else {
        n.setLastSavedLogId(n.getLastLogId())
    }
}

// 从物理化文件中恢复变量
func (n *Node) restoreDataFromFile() {
    path := n.getDataFilePath()
    if gfile.Exists(path) {
        bin := gfile.GetBinContents(path)
        if gCOMPRESS_SAVING {
            bin = gcompress.UnZlib(bin)
        }
        if bin != nil && len(bin) > 0 {
            glog.Println("restore data from", path)
            content := string(bin)
            var data = SaveInfo {
                LogList : make([]LogEntry, 0),
                Service : make(map[string]ServiceStruct),
                Peers   : make(map[string]interface{}),
                DataMap : make(map[string]string),
            }
            if gjson.DecodeTo(content, &data) == nil {
                n.setLastLogId(data.LastLogId)
                n.setLogCount(data.LogCount)
                n.setLastSavedLogId(data.LastLogId)
                n.setLastServiceLogId(data.LastServiceLogId)
                n.restoreLogList(&data)
                n.restoreService(&data)
                n.restoreDataMap(&data)
                n.restorePeer(&data)
            }
        }
    }
}

func (n *Node) restoreLogList(data *SaveInfo) {
    for _, v := range data.LogList {
        n.LogList.PushFront(v)
    }
}

func (n *Node) restoreService(data *SaveInfo) {
    for k, v := range data.Service {
        // 如果配置文件与数据文件有同一键名的配置，那么使用配置文件的Service配置
        if n.Service.Get(k) == nil {
            n.Service.Set(k, *n.serviceSructToService(&v))
        }
    }
}

func (n *Node) restorePeer(data *SaveInfo) {
    infoMap := make(map[string]NodeInfo)
    gjson.DecodeTo(gjson.Encode(data.Peers), &infoMap)
    for k, v := range infoMap {
        if !n.Peers.Contains(k) {
            n.Peers.Set(k, v)
        }
    }
}

func (n *Node) restoreDataMap(data *SaveInfo) {
    n.DataMap.BatchSet(data.DataMap)
}

// 使用logentry数组更新本地的日志列表
func (n *Node) updateFromLogEntriesJson(jsonContent string) error {
    array := make([]LogEntry, 0)
    err   := gjson.DecodeTo(jsonContent, &array)
    if err != nil {
        glog.Println(err)
        return err
    }
    if array != nil && len(array) > 0 {
        for _, v := range array {
            if v.Id > n.getLastLogId() {
                n.saveLogEntry(v)
            }
        }
    }
    return nil
}



