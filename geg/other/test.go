package main
import (
    "fmt"
)

type B struct {
    Name string
}

func (b B) Get() {
    fmt.Printf("b addr:%p\n", &b)
}
func main() {
    b := B{}
    b.Get()
    b.Get()

    return
    //view   := gmvc.NewView("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/mvc/view/user/")
    //tpl, _ := view.Template("info")
    //tpl.BindFunc("add", add)
    //fmt.Println(tpl.Parse(nil))
    ////t, err := template.New("text").Funcs(template.FuncMap{"add":add}).Parse(`{{add 1 2}}`)
    ////if err != nil {
    ////    panic(err)
    ////}
    ////t.Execute(os.Stdout, u)
}