// 数据同步需要注意的是：
// leader只有在通知完所有follower更新完数据之后，自身才会进行数据更新
// 因此leader
package gluster

import (
    "g/encoding/gjson"
    "time"
    "g/core/types/gmap"
    "g/os/gfile"
    "g/util/gtime"
    "g/os/glog"
)

// 保存日志数据
func (n *Node) saveLogEntry(entry LogEntry) {
    switch entry.Act {
        case gMSG_REPL_SET:
            glog.Println("setting log entry", entry)
            for k, v := range entry.Items.(map[string]interface{}) {
                n.KVMap.Set(k, v.(string))
            }

        case gMSG_REPL_REMOVE:
            glog.Println("removing log entry", entry)
            for _, v := range entry.Items.([]interface{}) {
                n.KVMap.Remove(v.(string))
            }

    }
    n.setLastLogId(entry.Id)
}

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
        LastServiceLogId    : n.getLastServiceLogId(),
        Service             : *n.Service.Clone(),
        Peers               : *n.Peers.Clone(),
        DataMap             : *n.KVMap.Clone(),
    }
    content := gjson.Encode(&data)
    gfile.PutContents(n.getDataFilePath(), *content)
    n.setLastSavedLogId(n.getLastLogId())
}

// 从物理化文件中恢复变量
func (n *Node) restoreDataFromFile() {
    path := n.getDataFilePath()
    if gfile.Exists(path) {
        content := gfile.GetContents(path)
        if content != nil {
            glog.Println("restore data from file:", path)
            var data = SaveInfo {
                Service : make(map[string]interface{}),
                Peers   : make(map[string]interface{}),
                DataMap : make(map[string]string),
            }
            content := string(content)
            if gjson.DecodeTo(&content, &data) == nil {
                n.setLastLogId(data.LastLogId)
                n.setLastSavedLogId(data.LastLogId)
                n.setLastServiceLogId(data.LastServiceLogId)
                n.restoreService(&data)
                n.restoreKVMap(&data)
                n.restorePeer(&data)
            }
        }
    }
}

func (n *Node) restoreService(data *SaveInfo) {
    serviceMap := gmap.NewStringInterfaceMap()
    servMap    := make(map[string]Service)
    gjson.DecodeTo(gjson.Encode(data.Service), &servMap)
    for k, v := range servMap {
        serviceMap.Set(k, v)
    }
    n.setService(serviceMap)
}

func (n *Node) restorePeer(data *SaveInfo) {
    peerMap := gmap.NewStringInterfaceMap()
    infoMap := make(map[string]NodeInfo)
    gjson.DecodeTo(gjson.Encode(data.Peers), &infoMap)
    for k, v := range infoMap {
        peerMap.Set(k, v)
    }
    n.setPeers(peerMap)
}

func (n *Node) restoreKVMap(data *SaveInfo) {
    dataMap := gmap.NewStringStringMap()
    dataMap.BatchSet(data.DataMap)
    n.setKVMap(dataMap)
}

// 使用logentry数组更新本地的日志列表
func (n *Node) updateFromLogEntriesJson(jsonContent *string) error {
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



