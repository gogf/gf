package ginstance


// 单例对象存储器
var instances = make(map[string]interface{})

// 获取单例对象
func Get(k string) interface{} {
    if v, ok := instances[k]; ok {
        return v
    } else {
        return nil
    }
}

// 设置单例对象
func Set(k string, v interface{}) {
    instances[k] = v
}