package main
import (
    "time"
    "fmt"
    "gitee.com/johng/gf/g/frame/gmvc"
    "gitee.com/johng/gf/g/net/gsession"
)
type User struct {
    Username, Password string
    RegTime time.Time
}
func add(i1, i2 int) int {
    return i1 + i2 + 1
}
func main() {
    fmt.Println(gsession.Id())
    view   := gmvc.NewView("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/mvc/view/user/")
    tpl, _ := view.Template("info")
    tpl.BindFunc("add", add)
    fmt.Println(tpl.Parse(nil))
    //t, err := template.New("text").Funcs(template.FuncMap{"add":add}).Parse(`{{add 1 2}}`)
    //if err != nil {
    //    panic(err)
    //}
    //t.Execute(os.Stdout, u)
}