package main
import (
    "fmt"
    "gitee.com/johng/gf/g/os/gview"
)

type B struct {

}

func add() int {return 1}

func main() {
    view   := gview.New("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/mvc/view/user/")
    tpl, _ := view.Template("info")
    tpl.BindFunc("include", add)
    fmt.Println(tpl.Parse(nil))
    //t, err := template.New("text").Funcs(template.FuncMap{"add":add}).Parse(`{{add 1 2}}`)
    //if err != nil {
    //    panic(err)
    //}
    //t.Execute(os.Stdout, u)
}