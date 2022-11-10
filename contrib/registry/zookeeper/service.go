package zookeeper

import (
	"encoding/json"
)

func unmarshal(data []byte) (c *Content, err error) {
	err = json.Unmarshal(data, &c)
	return
}

func marshal(c *Content) ([]byte, error) {
	return json.Marshal(c)
}
