package gconfig

// 配置对象
var config = make(map[string]interface{})

// 获取配置
func Get(k string) interface{} {
    if v, ok := config[k]; ok {
        return v
    } else {
        return nil
    }
}

// 设置配置
func Set(k string, v interface{}) {
    config[k] = v
}
