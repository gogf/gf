package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    s := `
[INFO] 2018-11-04 12:48:01  2018-11-04 12:48:01 ["handle eventMedlinker\\Med3Svr\\App\\Events\\Task\\InviteUserVerifyEvent::__set_state(array(\n   'deviceNum' => '25a715ff8d4835197299d7dab067841e',\n   'userIdList' => \n  array (\n    0 => '62506063',\n  ),\n   'dataId' => '',\n   'eventTime' => NULL,\n))"] [med3-svr-6c8c9b9f4f-fl6g9]

`
    t := gtime.ParseTimeFromContent(s)
    fmt.Println(t.String())
}