package gkvdb

import "time"

// 自动保存线程循环
func (db *DB) autoSavingLoop() {
    for {
        db.memt.sync()
        db.saveFileSpace()
        time.Sleep(gAUTO_SAVING_TIMEOUT*time.Millisecond)
    }
}