package user

import (
	"github.com/gogf/katyusha/krpc"
	"google.golang.org/grpc"
)

const (
	ServiceName = "demo"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(options ...grpc.DialOption) (*Client, error) {
	conn, err := krpc.Client.NewGrpcClientConn(ServiceName, options...)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) User() UserClient {
	return NewUserClient(c.conn)
}
