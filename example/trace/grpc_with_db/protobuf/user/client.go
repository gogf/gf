package user

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"google.golang.org/grpc"
)

const (
	ServiceName = "demo"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(options ...grpc.DialOption) (*Client, error) {
	conn, err := grpcx.Client.NewGrpcClientConn(ServiceName, options...)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) User() UserClient {
	return NewUserClient(c.conn)
}
