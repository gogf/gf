package gkvdb

import (
    "g/os/gfile"
    "g/encoding/gbinary"
)

// 元数据碎片
func (db *DB) addMtFileSpace(index int, size uint) {
    //defer db.saveFileSpace()
    db.mtsp.AddBlock(index, size)
}

func (db *DB) getMtFileSpace(size uint) (int, uint) {
    return db.mtsp.GetBlock(size)
}

// 数据碎片
func (db *DB) addDbFileSpace(index int, size uint) {
    //defer db.saveFileSpace()
    db.dbsp.AddBlock(index, size)
}

func (db *DB) getDbFileSpace(size uint) (int, uint) {
    return db.dbsp.GetBlock(size)
}

// 保存碎片数据到文件
func (db *DB) saveFileSpace() error {
    mtbuffer := db.mtsp.Export()
    dbbuffer := db.dbsp.Export()
    buffer   := make([]byte, 0)
    buffer    = append(buffer, gbinary.EncodeUint32(uint32(len(mtbuffer)))...)
    buffer    = append(buffer, gbinary.EncodeUint32(uint32(len(dbbuffer)))...)
    buffer    = append(buffer, mtbuffer...)
    buffer    = append(buffer, dbbuffer...)
    return gfile.PutBinContents(db.getSpaceFilePath(), buffer)
}

// 恢复碎片文件到内存
func (db *DB) restoreFileSpace() {
    buffer := gfile.GetBinContents(db.getSpaceFilePath())
    if len(buffer) > 8 {
        mtsize := gbinary.DecodeToUint32(buffer[0 : 4])
        dbsize := gbinary.DecodeToUint32(buffer[4 : 8])
        if mtsize > 0 {
            db.mtsp.Import(buffer[8 : 8 + mtsize])
        }
        if dbsize > 0 {
            db.mtsp.Import(buffer[8 + mtsize : 8 + mtsize + dbsize])
        }
    }
}