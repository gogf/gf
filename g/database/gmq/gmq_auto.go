package gmq

import (
    "time"
    "sync/atomic"
    "gitee.com/johng/gf/g/os/gfile"
)

// 自动清理过期的队列文件
func (mqg *MQGroup) startAutoClean() {
    go func() {
        minid := atomic.LoadUint64(&mqg.minid)
        for !mqg.isClosed() {
            id := atomic.LoadUint64(&mqg.minid)
            if id != minid {
                // 多个队列文件
                if id - minid >= gMQFILE_MAX_COUNT {
                    gfile.Remove(mqg.getFilePathById(minid))
                    minid = id
                }
                // 数据已全部使用完
                if id == atomic.LoadUint64(&mqg.maxid) {
                    gfile.Remove(mqg.getFilePathById(id))
                    minid = id
                }
            }
            time.Sleep(gMQFILE_AUTO_CLEAN_TIMEOUT*time.Second)
        }
    }()
}