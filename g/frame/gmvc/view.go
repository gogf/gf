package gmvc

import (
    "sync"
    "html/template"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/frame/gbase"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gconsole"
)

// 视图对象(一个请求一个视图对象，用完即销毁)
type View struct {
    gbase.Base
    mu   sync.RWMutex           // 并发互斥锁
    ctl  *Controller            // 所属控制器
    view *gview.View            // 底层视图对象
    data map[string]interface{} // 视图数据
}

// 创建一个MVC请求中使用的视图对象
func NewView(c *Controller) *View {
    // 视图目录路径查找优先级：配置文件参数viewpath、启动参数viewpath、当前程序运行目录
    path := gconsole.Option.Get("viewpath")
    if path == "" {
        path = gfile.SelfDir()
    }
    if r := c.Config.Get("viewpath"); r != nil {
        path = r.(string)
    }
    return &View{
        ctl  : c,
        view : gview.GetView(path),
        data : make(map[string]interface{}),
    }
}

// 批量绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assigns(data map[string]interface{}) {
    view.mu.Lock()
    defer view.mu.Unlock()
    for k, v := range data {
        view.data[k] = v
    }
}

// 绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assign(key string, value interface{}) {
    view.mu.Lock()
    defer view.mu.Unlock()
    view.data[key] = value
}

// 解析指定模板
func (view *View) Display(file string) error {
    // 查询模板
    tpl, err := view.view.Template(file)
    if err != nil {
        view.ctl.Response.WriteString("Tpl Parsing Error: " + err.Error())
        return err
    }
    // 绑定函数
    tpl.BindFunc("include", view.funcInclude)
    // 执行解析
    view.mu.RLock()
    defer view.mu.RUnlock()
    content, err := tpl.Parse(view.data)
    if err != nil {
        view.ctl.Response.WriteString("Tpl Parsing Error: " + err.Error())
        return err
    } else {
        view.ctl.Response.WriteString(content)
    }
    return nil
}

// 模板内置方法：include
func (view *View) funcInclude(file string) template.HTML {
    tpl, err := view.view.Template(file)
    if err != nil {
        return template.HTML(err.Error())
    }
    content, err := tpl.Parse(view.data)
    if err != nil {
        return template.HTML(err.Error())
    }
    return template.HTML(content)
}
