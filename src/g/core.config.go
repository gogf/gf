package g


type gConfig struct {
    Debug    bool
    Database struct {
        master struct{

        }
        slave struct{

        }
    }
}

// 框架设置项
var Config = gConfig {
    Debug : true,

}
