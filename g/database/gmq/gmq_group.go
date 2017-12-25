package gmq

import (
    "sync/atomic"
    "os"
    "gitee.com/johng/gf/g/os/gfile"
    "strconv"
    "gitee.com/johng/gf/g/os/gfilepool"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "strings"
    "bytes"
    "math"
    "sort"
)

// 根据需要写入的id获取对应的文件指针
func (mqg *MQGroup) getFileById(id uint64) (*os.File, error) {
    path  := mqg.getFilePathById(id)
    if file, err := gfilepool.OpenWithPool(path, os.O_RDWR|os.O_CREATE, gMQFILEPOOL_TIMEOUT); err == nil {
        return file.File(), nil
    } else {
        return nil, err
    }
}

// 根据id获取对应的队列文件的绝对路径
func (mqg *MQGroup) getFilePathById(id uint64) string {
    fnum  := int(id/uint64(gMQFILE_MAX_COUNT))
    path  := mqg.path + gfile.Separator + strconv.Itoa(fnum)
    return path
}

// 根据消息id计算索引在队列文件中的偏移量
func (mqg *MQGroup) getIndexOffsetById(id uint64) int64 {
    return int64(id%gMQFILE_MAX_COUNT)*gMQFILE_INDEX_ITEM_SIZE
}

// 初始化队列分类，获取当前队列分类的最小id及最大id
func (mqg *MQGroup) init() {
    minid   := uint64(0)
    maxid   := uint64(0)
    fnums   := make([]uint64, 0)
    files   := gfile.ScanDir(mqg.path)
    // 查找最大编号的队列文件
    for _, name := range files {
        if n, err := strconv.ParseUint(strings.Split(name, ".")[0], 10, 64); err == nil {
            fnums = append(fnums, n)
        }
    }
    sort.Slice(fnums, func(i, j int) bool { return fnums[i] < fnums[j] })
    if len(fnums) == 0 {
        return
    }
    // 查找当前队列文件的消息数量(最大id)
    if file, err := mqg.getFileById(uint64(fnums[len(fnums) - 1])*gMQFILE_MAX_COUNT); err == nil {
        if ixbuffer := gfile.GetBinContentByTwoOffsets(file, 0, gMQFILE_INDEX_LENGTH); ixbuffer != nil {
            zerob := make([]byte, gMQFILE_INDEX_ITEM_SIZE)
            for i := 0; i < gMQFILE_INDEX_LENGTH; i += gMQFILE_INDEX_ITEM_SIZE {
                maxid = uint64(fnums[len(fnums) - 1]*gMQFILE_MAX_COUNT)
                if bytes.Compare(zerob, ixbuffer[i : i + gMQFILE_INDEX_ITEM_SIZE]) == 0 {
                    if i > 0 {
                        maxid = uint64(fnums[len(fnums) - 1]*gMQFILE_MAX_COUNT) + uint64((i - gMQFILE_INDEX_ITEM_SIZE)/gMQFILE_INDEX_ITEM_SIZE)
                    }
                    if maxid != 0 {
                        break
                    }
                }

            }
        }
    }
    // 查找未使用的最小id
    for _, fnum := range fnums {
        if file, err := mqg.getFileById(uint64(fnum)*gMQFILE_MAX_COUNT); err == nil {
            if ixbuffer := gfile.GetBinContentByTwoOffsets(file, 0, gMQFILE_INDEX_LENGTH); ixbuffer != nil {
                for i := 0; i < gMQFILE_INDEX_LENGTH; i += gMQFILE_INDEX_ITEM_SIZE {
                    status := gbinary.DecodeToInt8(ixbuffer[i : i + 1])
                    if status == 1 {
                        minid = uint64(fnum*gMQFILE_MAX_COUNT) + uint64((i)/gMQFILE_INDEX_ITEM_SIZE)
                        if minid != 0 {
                            break
                        }
                    }
                }
            }
        }
        if minid != uint64(math.MaxUint64) {
            break
        }
    }
    mqg.minid = minid
    mqg.maxid = maxid
    mqg.startAutoClean()
}

func (mqg *MQGroup) getMinId() uint64 {
    return atomic.LoadUint64(&mqg.minid)
}

func (mqg *MQGroup) setMinId(id uint64) {
    atomic.StoreUint64(&mqg.minid, id)
}

func (mqg *MQGroup) getMaxId() uint64 {
    return atomic.LoadUint64(&mqg.maxid)
}

func (mqg *MQGroup) setMaxId(id uint64) {
    atomic.StoreUint64(&mqg.maxid, id)
}

// 获取队列的总数
func (mqg *MQGroup) Length() uint64 {
    if mqg.getMinId() > 0 {
        return mqg.getMaxId() - mqg.getMinId() + 1
    }
    return 0
}

// Push
func (mqg *MQGroup) Push(msg []byte) error {
    // 锁定队列文件
    mqg.mu.Lock()
    defer mqg.mu.Unlock()

    id        := mqg.maxid + 1
    file, err := mqg.getFileById(id)
    if err != nil {
        return err
    }
    ixoffset := mqg.getIndexOffsetById(id)

    // 消息数据写到文件末尾
    dataOffset, err := file.Seek(0, 2)
    if err != nil {
        return err
    }
    if dataOffset < gMQFILE_INDEX_LENGTH {
        dataOffset = gMQFILE_INDEX_LENGTH
        // 文件需要初始化，索引域初始化为0
        if _, err := file.WriteAt(make([]byte, gMQFILE_INDEX_LENGTH), 0); err != nil {
            return err
        }
    }
    if _, err := file.WriteAt(msg, dataOffset); err != nil {
        return err
    }
    // 数据写入成功后再写入索引域
    bits   := make([]gbinary.Bit, 0)
    bits    = gbinary.EncodeBits(bits, uint(dataOffset),  40)
    bits    = gbinary.EncodeBits(bits, uint(len(msg)),    24)
    indexb := append(gbinary.EncodeInt8(1), gbinary.EncodeBitsToBytes(bits)...)
    if _, err := file.WriteAt(indexb, ixoffset); err != nil {
        return err
    }
    // 执行成功之后最大id才会递增
    mqg.maxid++
    return nil
}

// POP，队列没有数据返回nil
func (mqg *MQGroup) Pop() []byte {
    mqg.mu.Lock()
    defer mqg.mu.Unlock()
    // 没有数据可获取
    if mqg.minid == 0 || mqg.minid == mqg.maxid {
        return nil
    }
    if buffer := mqg.get(mqg.minid); buffer != nil {
        if mqg.remove(mqg.minid) == nil {
            mqg.minid++
            return buffer
        }
    }
    return nil
}

// 根据消息id查询消息
func (mqg *MQGroup) get(id uint64) []byte {
    file, err := mqg.getFileById(id)
    if err != nil {
        return nil
    }
    offset := mqg.getIndexOffsetById(id)
    if ixbuffer := gfile.GetBinContentByTwoOffsets(file, offset, offset + gMQFILE_INDEX_ITEM_SIZE); ixbuffer != nil {
        status := gbinary.DecodeToInt8(ixbuffer[0 : 1])
        if status > 0 {
            bits   := gbinary.DecodeBytesToBits(ixbuffer[1 : ])
            start  := gbinary.DecodeBits(bits[0 : 40])
            size   := gbinary.DecodeBits(bits[40 : ])
            end    := start + size
            return gfile.GetBinContentByTwoOffsets(file, int64(start), int64(end))
        }
    }
    return nil
}

// 根据消息id删除消息
func (mqg *MQGroup) remove(id uint64) error {
    file, err := mqg.getFileById(id)
    if err != nil {
        return err
    }
    // 标记对应索引字段为删除(软删除)
    if _, err := file.WriteAt([]byte{byte(0)}, mqg.getIndexOffsetById(id)); err != nil {
        return err
    }
    return nil
}

// 关闭队列分类，自动回收资源
func (mqg *MQGroup) Close() {
    mqg.mu.Lock()
    defer mqg.mu.Unlock()
    mqg.closed = true
}

// 队列分类是否已关闭
func (mqg *MQGroup) isClosed() bool {
    mqg.mu.RLock()
    defer mqg.mu.RUnlock()
    return mqg.closed
}