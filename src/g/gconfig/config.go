package gconfig

// 全局配置管理对象
var Config gConfig

type gConfig struct {
    config map[string]interface{}
}

// 获取配置
func (i gConfig) Get(k string) interface{} {
    if v, ok := i.config[k]; ok {
        return v
    } else {
        return nil
    }
}

// 设置配置
func (i gConfig) Set(k string, v interface{}) {
    i.config[k] = v
}
