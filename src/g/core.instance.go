package g

// 全局单例操作对象
var Instance gInstance

// 单例对象类型
type gInstance struct {
    instances map[string]interface{}
}

// 获取单例对象
func (i gInstance) Get(k string) interface{} {
    if v, ok := i.instances[k]; ok {
        return v
    } else {
        return nil
    }
}

// 设置单例对象
func (i gInstance) Set(k string, v interface{}) {
    i.instances[k] = v
}