package gf


/*
    自定义类型命名规范，使用前缀表示类型，所有gf框架定义的类型都需要使用g作为前缀，并随后接类型缩写词：
    gstXxx   结构体
    gitXxx   接口
    gmapXxx  哈希表
    gfuncXxx 函数
 */
// 默认空函数类型，一般用于入口函数(特别是控制器入口函数)
type GfuncEmpty func()
