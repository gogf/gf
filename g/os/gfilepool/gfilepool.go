package gfilepool

import (
    "os"
    "time"
    "strconv"
    "gitee.com/johng/gf/g/util/gtime"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/glist"
)

// 文件指针池
type Pool struct {
    path    string          // 文件绝对路径
    flag    int             // 文件打开标识
    list    *glist.SafeList // 可用/闲置的文件指针链表
    idlemax int             // 闲置最大时间，超过该时间则被系统回收(秒)
    closed  bool            // 连接池是否已关闭
}

// 文件指针池指针
type File struct {
    pool   *Pool     // 所属池
    file   *os.File  // 指针对象
    expire int64     // 过期时间
}

// 全局指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
var pools *gmap.StringInterfaceMap = gmap.NewStringInterfaceMap()

// 获得文件对象，并自动创建指针池
func OpenWithPool(path string, flag int, expire int) (*File, error) {
    key    := path + strconv.Itoa(flag) + strconv.Itoa(expire)
    result := pools.Get(key)
    if result != nil {
        return result.(*Pool).File()
    }
    pool := New(path, flag, expire)
    pools.Set(key, pool)
    return pool.File()
}

// 创建一个文件指针池，expire < 0表示不过期，expire = 0表示使用完立即回收，expire > 0表示超时回收
func New(path string, flag int, expire int) *Pool {
    r := &Pool {
        path    : path,
        flag    : flag,
        list    : glist.NewSafeList(),
        idlemax : expire,
    }
    // 独立的线程执行过期清理工作
    if expire != -1 {
        go func(p *Pool) {
            for !p.closed {
                r := p.list.Front()
                if r != nil && r.Value != nil {
                    f := r.Value.(*File)
                    if f.expire <= gtime.Second() {
                        if f.file != nil {
                            f.file.Close()
                        }
                        p.list.Remove(r)
                        continue
                    }
                }
                time.Sleep(3 * time.Second)
            }
        }(r)
    }
    return r
}

// 获得一个文件打开指针
func (p *Pool) File() (*File, error) {
    if p.list.Len() > 0 {
        for {
            r := p.list.PopBack()
            if r != nil {
                f := r.(*File)
                if f.expire > gtime.Second() {
                    return f, nil
                } else if f.file != nil {
                    f.file.Close()
                    f.file = nil
                }
            } else {
                break;
            }
        }
    }
    file, err := os.OpenFile(p.path, p.flag, 0755)
    if err != nil {
        return nil, err
    }
    return &File {
        pool : p,
        file : file,
    }, nil
}

// 关闭指针池
func (p *Pool) Close() {
    p.closed = true
}

// 获得底层文件指针
func (f *File) File() *os.File {
    return f.file
}

// 关闭指针链接(软关闭)
func (f *File) Close() {
    f.expire = gtime.Second() + int64(f.pool.idlemax)
    f.pool.list.PushFront(f)
}