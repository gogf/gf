package funcs

import (
	"encoding/json"
	"fmt"

	"github.com/jin502437344/gf/.example/net/gtcp/pkg_operations/common/types"
	"github.com/jin502437344/gf/net/gtcp"
)

// 自定义格式发送消息包
func SendPkg(conn *gtcp.Conn, act string, data ...string) error {
	s := ""
	if len(data) > 0 {
		s = data[0]
	}
	msg, err := json.Marshal(types.Msg{
		Act:  act,
		Data: s,
	})
	if err != nil {
		panic(err)
	}
	return conn.SendPkg(msg)
}

// 自定义格式接收消息包
func RecvPkg(conn *gtcp.Conn) (msg *types.Msg, err error) {
	if data, err := conn.RecvPkg(); err != nil {
		return nil, err
	} else {
		msg = &types.Msg{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			return nil, fmt.Errorf("invalid package structure: %s", err.Error())
		}
		return msg, err
	}
}
